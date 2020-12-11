package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.f110.dev/heimdallr/pkg/cert"
	"go.f110.dev/heimdallr/pkg/config/configv2"
	"go.f110.dev/heimdallr/pkg/database"
	"go.f110.dev/heimdallr/pkg/database/memory"
	"go.f110.dev/heimdallr/pkg/rpc"
	"go.f110.dev/heimdallr/pkg/rpc/rpcclient"
	"go.f110.dev/heimdallr/pkg/session"
)

func TestInit(t *testing.T) {
	Init(
		&configv2.Config{AccessProxy: &configv2.AccessProxy{}},
		session.NewSecureCookieStore([]byte(""), []byte(""), ""),
		memory.NewUserDatabase(),
		memory.NewTokenDatabase(),
		&testRevokedCertClient{},
	)
}

func TestAuthenticate(t *testing.T) {
	s := session.NewSecureCookieStore([]byte("test"), []byte("testtesttesttesttesttesttesttest"), "example.com")
	u := memory.NewUserDatabase()
	rc := &testRevokedCertClient{}
	caCert, caPrivateKey, err := cert.CreateCertificateAuthority("for test", "test", "", "jp")
	if err != nil {
		t.Fatal(err)
	}
	cp := x509.NewCertPool()
	cp.AddCert(caCert)
	conf := &configv2.Config{
		AccessProxy: &configv2.AccessProxy{
			ServerNameHost: "proxy.example.com",
			Backends: []*configv2.Backend{
				{
					Name: "test",
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
						{Name: "ok_but_nobind", Locations: []configv2.Location{{Get: "/no_bind"}}},
					},
				},
				{
					Name:               "topsecret",
					MaxSessionDuration: &configv2.Duration{Duration: 1 * time.Minute},
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
					},
				},
				{
					Name:          "root",
					AllowRootUser: true,
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
					},
				},
				{
					Name:         "public",
					DisableAuthn: true,
				},
				{
					Name:         "public-path",
					DisableAuthn: true,
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
					},
				},
			},
		},
		AuthorizationEngine: &configv2.AuthorizationEngine{
			RootUsers: []string{"root@example.com"},
			Roles: []*configv2.Role{
				{
					Name: "test",
					Bindings: []*configv2.Binding{
						{Backend: "test", Permission: "ok"},
						{Backend: "topsecret", Permission: "ok"},
					},
				},
			},
		},
		CertificateAuthority: &configv2.CertificateAuthority{
			Local: &configv2.CertificateAuthorityLocal{
				Certificate: caCert,
				PrivateKey:  caPrivateKey,
				CertPool:    cp,
			},
		},
	}
	err = conf.AccessProxy.Setup(conf.AccessProxy.Backends)
	require.NoError(t, err)
	err = conf.AuthorizationEngine.Setup(conf.AuthorizationEngine.Roles, []*configv2.RPCPermission{})
	require.NoError(t, err)
	Init(conf, s, u, nil, rc)
	_ = u.Set(nil, &database.User{Id: "foobar@example.com", Roles: []string{"test", "unknown"}})

	t.Run("DisableAuthentication", func(t *testing.T) {
		t.Parallel()

		t.Run("any path", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://public.proxy.example.com/ok", nil)
			_, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("with restricted path", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://public-path.proxy.example.com/ok", nil)
			_, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("with not allowed path", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://public-path.proxy.example.com/disallow", nil)
			_, _, err := Authenticate(context.TODO(), req)
			if err != ErrNotAllowed {
				t.Fatalf("expected ErrNotAllowed: %v", err)
			}
		})
	})

	t.Run("Cookie", func(t *testing.T) {
		t.Parallel()

		t.Run("normal", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			c, err := s.Cookie(session.New("foobar@example.com"))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(c)

			user, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
			if user.Id != "foobar@example.com" {
				t.Errorf("expect foobar@example.com: %s", user.Id)
			}

			user, _, err = Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
			if user.Id != "foobar@example.com" {
				t.Errorf("expect foobar@example.com: %s", user.Id)
			}
		})

		t.Run("normal with root user", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://root.proxy.example.com/ok", nil)
			c, err := s.Cookie(session.New("root@example.com"))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(c)
			user, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
			if user.Id != "root@example.com" {
				t.Errorf("expect foobar@example.com: %s", user.Id)
			}
			if !user.RootUser {
				t.Error("expect root@example.com is root")
			}
		})

		t.Run("not have cookie", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)

			_, _, err := Authenticate(context.TODO(), req)
			if err != ErrSessionNotFound {
				t.Errorf("expect session not found: %v", err)
			}
		})

		t.Run("user not found", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			c, err := s.Cookie(session.New("foo@example.com"))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(c)
			_, _, err = Authenticate(context.TODO(), req)
			if err != ErrUserNotFound {
				t.Errorf("expect user not found: %v", err)
			}
		})

		t.Run("unknown host", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://unknownhost.example.com/ok", nil)

			_, _, err := Authenticate(context.TODO(), req)
			if err != ErrHostnameNotFound {
				t.Errorf("expect hostname not found: %v", err)
			}
		})
	})

	t.Run("client certificate auth", func(t *testing.T) {
		t.Parallel()

		newClientCert := func(t *testing.T, id string) *x509.Certificate {
			subject := pkix.Name{CommonName: id}
			pemEncodedCSRBytes, _, err := cert.CreateCertificateRequest(subject, []string{})
			if err != nil {
				t.Fatal(err)
			}
			block, _ := pem.Decode(pemEncodedCSRBytes)
			csr, err := x509.ParseCertificateRequest(block.Bytes)
			if err != nil {
				t.Fatal(err)
			}
			clientCert, err := cert.SigningCertificateRequest(csr, conf.CertificateAuthority.Local)
			if err != nil {
				t.Fatal(err)
			}

			return clientCert
		}

		t.Run("normal", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.TLS = &tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{newClientCert(t, "foobar@example.com")},
			}
			user, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
			if user.Id != "foobar@example.com" {
				t.Errorf("expect foobar@example.com: %s", user.Id)
			}
		})

		t.Run("normal with revoked certificate", func(t *testing.T) {
			t.Parallel()

			clientCert := newClientCert(t, "foobar@example.com")
			rc.revokedCert = append(rc.revokedCert, &rpcclient.RevokedCert{SerialNumber: clientCert.SerialNumber})

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.TLS = &tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{clientCert},
			}
			_, _, err = Authenticate(context.TODO(), req)
			if err == nil {
				t.Fatal("expected to occurred error but not")
			}
			if err != ErrInvalidCertificate {
				t.Fatalf("expected ErrInvalidCertificate: %v", err)
			}
		})

		t.Run("access to Secure backend with Cookie", func(t *testing.T) {
			t.Parallel()

			clientCert := newClientCert(t, "foobar@example.com")
			req := httptest.NewRequest(http.MethodGet, "http://topsecret.proxy.example.com/ok", nil)
			req.TLS = &tls.ConnectionState{
				PeerCertificates: []*x509.Certificate{clientCert},
			}
			c, err := s.Cookie(session.New("foobar@example.com"))
			if err != nil {
				t.Fatal(err)
			}
			req.AddCookie(c)

			_, _, err = Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Authorization header", func(t *testing.T) {
		t.Parallel()

		t.Run("normal", func(t *testing.T) {
			t.Parallel()

			err := u.SetAccessToken(context.Background(), &database.AccessToken{Value: t.Name(), UserId: "foobar@example.com"})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.Header.Set("Authorization", "LP-TOKEN")
			req.Header.Set("X-LP-TOKEN", t.Name())
			user, _, err := Authenticate(context.TODO(), req)
			if err != nil {
				t.Fatal(err)
			}
			if user.Id != "foobar@example.com" {
				t.Errorf("expect foobar@example.com: %s", user.Id)
			}
		})

		t.Run("header not found", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.Header.Set("Authorization", "LP-TOKEN")
			_, _, err = Authenticate(context.TODO(), req)
			if err == nil {
				t.Fatal("expected to occurred error but not")
			}
			if err != ErrUserNotFound {
				t.Fatalf("expected ErrUserNotFound: %v", err)
			}
		})

		t.Run("invalid token", func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.Header.Set("Authorization", "LP-TOKEN")
			req.Header.Set("X-LP-Token", "unknown-token")
			_, _, err = Authenticate(context.TODO(), req)
			if err == nil {
				t.Fatal("expected to occurred error but not")
			}
			if err != ErrUserNotFound {
				t.Fatalf("expected ErrUserNotFound: %v", err)
			}
		})

		t.Run("user not found", func(t *testing.T) {
			t.Parallel()

			err := u.SetAccessToken(context.Background(), &database.AccessToken{Value: "dummy-token", UserId: "piyo@example.com"})
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodGet, "http://test.proxy.example.com/ok", nil)
			req.Header.Set("Authorization", "LP-TOKEN")
			req.Header.Set("X-LP-TOKEN", "dummy-token")
			_, _, err = Authenticate(context.TODO(), req)
			if err == nil {
				t.Fatal("expected to occurred error but not")
			}
			if err != ErrUserNotFound {
				t.Fatalf("expected ErrUserNotFound: %v", err)
			}
		})
	})
}

func TestAuthenticator_AuthenticateForSocket(t *testing.T) {
	s := session.NewSecureCookieStore([]byte("test"), []byte("testtesttesttesttesttesttesttest"), "example.com")
	u := memory.NewUserDatabase()
	token := memory.NewTokenDatabase()
	conf := &configv2.Config{
		AccessProxy: &configv2.AccessProxy{
			ServerNameHost: "proxy.example.com",
			Backends: []*configv2.Backend{
				{
					Name:   "test",
					Socket: true,
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
						{Name: "ok_but_nobind", Locations: []configv2.Location{{Get: "/no_bind"}}},
					},
				},
			},
		},
		AuthorizationEngine: &configv2.AuthorizationEngine{
			Roles: []*configv2.Role{
				{
					Name: "test",
					Bindings: []*configv2.Binding{
						{Backend: "test"},
					},
				},
			},
		},
	}
	err := conf.AccessProxy.Setup(conf.AccessProxy.Backends)
	require.NoError(t, err)
	err = conf.AuthorizationEngine.Setup(conf.AuthorizationEngine.Roles, []*configv2.RPCPermission{})
	require.NoError(t, err)

	Init(conf, s, u, token, nil)
	err = u.Set(nil, &database.User{Id: "foobar@example.com", Roles: []string{"test", "unknown"}})
	require.NoError(t, err)
	err = u.Set(nil, &database.User{Id: "piyo@example.com", Roles: []string{}})
	require.NoError(t, err)

	t.Run("empty token", func(t *testing.T) {
		t.Parallel()

		_, _, err := AuthenticateSocket(context.Background(), "", "")
		require.Error(t, err)
		require.Equal(t, ErrInvalidToken, err)
	})

	t.Run("empty host", func(t *testing.T) {
		t.Parallel()

		_, _, err := AuthenticateSocket(context.Background(), "dummy", "")
		require.Error(t, err)
		require.Equal(t, ErrHostnameNotFound, err)
	})

	t.Run("unknown host", func(t *testing.T) {
		t.Parallel()

		_, _, err := AuthenticateSocket(context.Background(), "dummy", "unknown")
		require.Error(t, err)
		require.Equal(t, ErrHostnameNotFound, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()

		_, _, err := AuthenticateSocket(context.Background(), "dummy", "test.proxy.example.com")
		require.Error(t, err)
		require.Equal(t, ErrInvalidToken, err)
	})

	t.Run("Unknown user", func(t *testing.T) {
		t.Parallel()

		newToken, err := token.IssueToken(context.Background(), "", "")
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = AuthenticateSocket(context.Background(), newToken.Token, "test.proxy.example.com")
		require.Error(t, err)
		require.Equal(t, ErrUserNotFound, err)
	})

	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		newCode, err := token.NewCode(context.Background(), "foobar@example.com", "", "")
		if err != nil {
			t.Fatal(err)
		}
		newToken, err := token.IssueToken(context.Background(), newCode.Code, "")
		if err != nil {
			t.Fatal(err)
		}
		backend, user, err := AuthenticateSocket(context.Background(), newToken.Token, "test.proxy.example.com")
		require.NoError(t, err)
		assert.Equal(t, "test", backend.Name)
		assert.Equal(t, "foobar@example.com", user.Id)
	})
}

func TestAuthInterceptor_UnaryInterceptor(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	u := memory.NewUserDatabase(database.SystemUser)
	token := memory.NewTokenDatabase()
	conf := &configv2.Config{
		AccessProxy: &configv2.AccessProxy{
			ServerNameHost: "proxy.example.com",
			Credential: &configv2.Credential{
				SigningPrivateKey: privateKey,
				InternalToken:     "rpc-internal-token",
				SigningPublicKey:  privateKey.PublicKey,
			},
			Backends: []*configv2.Backend{
				{
					Name: "test",
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
						{Name: "ok_but_nobind", Locations: []configv2.Location{{Get: "/no_bind"}}},
					},
				},
			},
		},
		AuthorizationEngine: &configv2.AuthorizationEngine{
			Roles: []*configv2.Role{
				{
					Name: "test",
					Bindings: []*configv2.Binding{
						{RPC: "test"},
					},
				},
			},
			RPCPermissions: []*configv2.RPCPermission{
				{
					Name:  "test",
					Allow: []string{"test"},
				},
			},
		},
	}
	err = conf.AccessProxy.Setup(conf.AccessProxy.Backends)
	require.NoError(t, err)
	err = conf.AuthorizationEngine.Setup(conf.AuthorizationEngine.Roles, conf.AuthorizationEngine.RPCPermissions)
	require.NoError(t, err)
	Init(conf, nil, u, token, nil)
	_ = u.Set(nil, &database.User{Id: "foobar@example.com", Roles: []string{"test"}})

	okHandler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return true, nil
	}

	t.Run("with access token", func(t *testing.T) {
		t.Parallel()

		newCode, err := token.NewCode(context.Background(), "foobar@example.com", "", "")
		if err != nil {
			t.Fatal(err)
		}
		newToken, err := token.IssueToken(context.Background(), newCode.Code, "")
		if err != nil {
			t.Fatal(err)
		}

		md := metadata.New(map[string]string{rpc.TokenMetadataKey: newToken.Token})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		v, err := UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}
		res, ok := v.(bool)
		if !ok {
			t.Fatal("response should be bool")
		}
		if !res {
			t.Fatal("unexpected response")
		}

		v, err = UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}
		res, ok = v.(bool)
		if !ok {
			t.Fatal("response should be bool")
		}
		if !res {
			t.Fatal("unexpected response")
		}
	})

	t.Run("with jwt", func(t *testing.T) {
		t.Parallel()

		claim := jwt.NewWithClaims(jwt.SigningMethodES256, &jwt.StandardClaims{
			Id:        "foobar@example.com",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
		})
		jwtToken, err := claim.SignedString(conf.AccessProxy.Credential.SigningPrivateKey)
		if err != nil {
			t.Fatal(err)
		}

		md := metadata.New(map[string]string{rpc.JwtTokenMetadataKey: jwtToken})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		v, err := UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}
		res, ok := v.(bool)
		if !ok {
			t.Fatal("response should be bool")
		}
		if !res {
			t.Fatal("unexpected response")
		}
	})

	t.Run("with internal token", func(t *testing.T) {
		t.Parallel()

		md := metadata.New(map[string]string{rpc.InternalTokenMetadataKey: conf.AccessProxy.Credential.InternalToken})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		v, err := UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/proxy.rpc.certificateauthority.watchrevokedcert"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}
		res, ok := v.(bool)
		if !ok {
			t.Fatal("response should be bool")
		}
		if !res {
			t.Fatal("unexpected response")
		}
	})

	t.Run("not provide metadata", func(t *testing.T) {
		t.Parallel()

		_, err := UnaryInterceptor(context.Background(), nil, nil, nil)
		if err == nil {
			t.Fatal("expect to occurred error but not")
		}
		if !errors.Is(err, unauthorizedError.Err()) {
			t.Fatalf("expect unauthorizedError: %v", err)
		}
	})

	t.Run("not provide token", func(t *testing.T) {
		t.Parallel()

		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, nil)
		if err == nil {
			t.Fatal("expect to occurred error but not")
		}
		if !errors.Is(err, unauthorizedError.Err()) {
			t.Fatalf("expect unauthorizedError: %v", err)
		}
	})

	t.Run("health check methods should not check a clearance", func(t *testing.T) {
		t.Parallel()

		methods := []string{"/grpc.health.v1.Health/Check", "/proxy.rpc.Admin/Ping"}

		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		for _, m := range methods {
			v, err := UnaryInterceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: m}, func(_ context.Context, _ interface{}) (interface{}, error) {
				return true, nil
			})
			if err != nil {
				t.Fatal(err)
			}
			res, ok := v.(bool)
			if !ok {
				t.Fatal("response should be bool")
			}
			if !res {
				t.Fatal("unexpected response")
			}
		}
	})
}

func TestAuthInterceptor_StreamInterceptor(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	u := memory.NewUserDatabase(database.SystemUser)
	token := memory.NewTokenDatabase()
	conf := &configv2.Config{
		AccessProxy: &configv2.AccessProxy{
			ServerNameHost: "proxy.example.com",
			Credential: &configv2.Credential{
				SigningPrivateKey: privateKey,
				InternalToken:     "rpc-internal-token",
			},
			Backends: []*configv2.Backend{
				{
					Name: "test",
					Permissions: []*configv2.Permission{
						{Name: "ok", Locations: []configv2.Location{{Get: "/ok"}}},
						{Name: "ok_but_nobind", Locations: []configv2.Location{{Get: "/no_bind"}}},
					},
				},
			},
		},
		AuthorizationEngine: &configv2.AuthorizationEngine{
			Roles: []*configv2.Role{
				{
					Name: "test",
					Bindings: []*configv2.Binding{
						{RPC: "test"},
					},
				},
			},
			RPCPermissions: []*configv2.RPCPermission{
				{
					Name:  "test",
					Allow: []string{"test"},
				},
			},
		},
	}
	err = conf.AccessProxy.Setup(conf.AccessProxy.Backends)
	require.NoError(t, err)
	err = conf.AuthorizationEngine.Setup(conf.AuthorizationEngine.Roles, conf.AuthorizationEngine.RPCPermissions)
	require.NoError(t, err)
	Init(conf, nil, u, token, nil)

	okHandler := func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	}

	t.Run("with internal token", func(t *testing.T) {
		t.Parallel()

		md := metadata.New(map[string]string{rpc.InternalTokenMetadataKey: conf.AccessProxy.Credential.InternalToken})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		err := StreamInterceptor(nil, &testServerStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/test"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}

		err = StreamInterceptor(nil, &testServerStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/test"}, okHandler)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("not provide metadata", func(t *testing.T) {
		t.Parallel()

		err := StreamInterceptor(nil, &testServerStream{ctx: context.Background()}, nil, nil)
		if err == nil {
			t.Fatal("expect to occurred error but not")
		}
		if !errors.Is(err, unauthorizedError.Err()) {
			t.Fatalf("expect unauthorizedError: %v", err)
		}
	})

	t.Run("not provide token", func(t *testing.T) {
		t.Parallel()

		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		err := StreamInterceptor(nil, &testServerStream{ctx: ctx}, &grpc.StreamServerInfo{FullMethod: "/test"}, nil)
		if err == nil {
			t.Fatal("expect to occurred error but not")
		}
		if !errors.Is(err, unauthorizedError.Err()) {
			t.Fatalf("expect unauthorizedError: %v", err)
		}
	})
}
