package scenario

import (
	"fmt"
	"net/http"
	"testing"

	"go.f110.dev/heimdallr/e2e/framework"
	"go.f110.dev/heimdallr/pkg/config/configv2"
	"go.f110.dev/heimdallr/pkg/database"
	"go.f110.dev/heimdallr/pkg/session"
)

func TestL7ReverseProxy(t *testing.T) {
	f := framework.New(t)
	defer f.Execute()

	f.Describe("L7 Reverse Proxy", func(s *framework.Scenario) {
		s.BeforeEach(func(m *framework.Matcher) bool { return m.Must(f.Proxy.Reload()) })
		s.Defer(func() { f.Proxy.Cleanup() })

		s.Context("authorization flow", func(s *framework.Scenario) {
			s.BeforeAll(func(m *framework.Matcher) bool {
				f.Proxy.Backend(&configv2.Backend{Name: "test", Permissions: []*configv2.Permission{{Name: "all", Locations: []configv2.Location{{Any: "/"}}}}})
				f.Proxy.Role(&configv2.Role{Name: "test", Bindings: []*configv2.Binding{{Backend: "test", Permission: "all"}}})
				f.Proxy.User(&database.User{Id: "test@f110.dev", Roles: []string{"test"}})
				return true
			})
			s.AfterAll(func(m *framework.Matcher) bool { return f.Proxy.ClearConf() })

			s.Step("request backend url", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return f.Agents.User("test@f110.dev").Get(m, f.Proxy.URL("test"))
				})

				s.It("should redirect to entry point", func(m *framework.Matcher) bool {
					m.StatusCode(http.StatusSeeOther)
					u, err := m.LastResponse().Location()
					m.NoError(err)
					return m.Contains(u.String(), f.Proxy.URL("", "/auth"))
				})
			})

			s.Step("enter authorization flow", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					m.Must(f.Agents.User("test@f110.dev").FollowRedirect(m))
					f.Agents.User("test@f110.dev").SaveCookie()
					return true
				})

				s.It("should redirect to OpenID Connect auth endpoint", func(m *framework.Matcher) bool {
					m.StatusCode(http.StatusFound)
					u, err := m.LastResponse().Location()
					m.NoError(err)
					return m.Contains(u.String(), "custom-idp/auth")
				})

				s.It("receive the cookie", func(m *framework.Matcher) bool {
					cookie := m.LastResponse().FindCookie(session.CookieName)
					m.NotNil(cookie)
					m.Equal(f.Proxy.DomainHost, cookie.Domain)
					m.True(cookie.HttpOnly)
					m.True(cookie.Secure)
					sess, err := f.Agents.DecodeCookieValue(cookie.Name, cookie.Value)
					m.NoError(err)
					m.Empty(sess.Id)
					return m.NotEmpty(sess.Unique)
				})
			})

			s.Step("show authorization view of identity provider", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return m.Must(f.Agents.User("test@f110.dev").FollowRedirect(m))
				})

				s.It("should get a page", func(m *framework.Matcher) bool {
					return m.StatusCode(http.StatusOK)
				})
			})

			s.Step("login identity provider", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					agent := f.Agents.User("test@f110.dev")
					authResponse := &framework.AuthResponse{}
					m.Must(agent.ParseLastResponseBody(authResponse))
					return agent.Post(m, authResponse.LoginURL, fmt.Sprintf(`{"id":"test@f110.dev","query":"%s"}`, authResponse.Query))
				})

				s.It("should redirect to callback", func(m *framework.Matcher) bool {
					m.StatusCode(http.StatusFound)
					u, err := m.LastResponse().Location()
					m.NoError(err)
					return m.Contains(u.String(), "auth/callback")
				})
			})

			s.Step("follow redirect", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return m.Must(f.Agents.User("test@f110.dev").FollowRedirect(m))
				})

				s.It("should redirect to backend", func(m *framework.Matcher) bool {
					m.StatusCode(http.StatusFound)
					u, err := m.LastResponse().Location()
					m.NoError(err)
					return m.Equal(f.Proxy.URL("test"), u.String())
				})
			})
		})

		s.Context("access to the unknown backend", func(s *framework.Scenario) {
			s.Context("by authenticated user", func(s *framework.Scenario) {
				s.BeforeAll(func(m *framework.Matcher) bool {
					f.Proxy.Backend(&configv2.Backend{Name: "test", Permissions: []*configv2.Permission{{Name: "all", Locations: []configv2.Location{{Any: "/"}}}}})
					f.Proxy.Role(&configv2.Role{Name: "test", Bindings: []*configv2.Binding{{Backend: "test", Permission: "all"}}})
					f.Proxy.User(&database.User{Id: "test@f110.dev", Roles: []string{"test"}})
					return true
				})
				s.AfterAll(func(m *framework.Matcher) bool { return f.Proxy.ClearConf() })

				s.It("should close connection", func(m *framework.Matcher) bool {
					return f.Agents.Authorized("test@f110.dev").Get(m, f.Proxy.URL("unknown"))
				})
			})

			s.Context("by unauthenticated agent", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return f.Agents.Unauthorized().Get(m, f.Proxy.URL("unknown"))
				})

				s.It("should close connection", func(m *framework.Matcher) bool {
					return m.ResetConnection()
				})
			})
		})

		s.Context("access to the backend", func(s *framework.Scenario) {
			s.BeforeAll(func(m *framework.Matcher) bool {
				f.Proxy.Backend(&configv2.Backend{
					Name: "test",
					HTTP: []*configv2.HTTPBackend{{Path: "/"}},
					Permissions: []*configv2.Permission{
						{Name: "all", Locations: []configv2.Location{{Any: "/"}}},
					},
				})
				f.Proxy.Role(&configv2.Role{Name: "test", Bindings: []*configv2.Binding{{Backend: "test", Permission: "all"}}})
				f.Proxy.Role(&configv2.Role{Name: "test2"})
				f.Proxy.User(&database.User{Id: "test@f110.dev", Roles: []string{"test"}})
				f.Proxy.User(&database.User{Id: "test2@f110.dev", Roles: []string{"test2"}})
				return true
			})
			s.AfterAll(func(m *framework.Matcher) bool { return f.Proxy.ClearConf() })

			s.Context("by unauthenticated agent", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return f.Agents.Unauthorized().Get(m, f.Proxy.URL("test"))
				})

				s.It("should redirect to IdP", func(m *framework.Matcher) bool {
					m.StatusCode(http.StatusSeeOther)
					u, err := m.LastResponse().Location()
					m.NoError(err)
					return m.Equal(f.Proxy.URL("test"), u.Query().Get("from"))
				})
			})

			s.Context("by authenticated user", func(s *framework.Scenario) {
				s.Context("who allowed an access", func(s *framework.Scenario) {
					s.Subject(func(m *framework.Matcher) bool {
						return f.Agents.Authorized("test@f110.dev").Get(m, f.Proxy.URL("test", "index.html"))
					})

					s.It("should proxy to backend", func(m *framework.Matcher) bool {
						return m.Equal(http.StatusBadGateway, m.LastResponse().StatusCode, "returns status 502 (BadGateway) because the upstream is down")
					})
				})

				s.Context("who not allowed an access", func(s *framework.Scenario) {
					s.Subject(func(m *framework.Matcher) bool {
						return f.Agents.Authorized("test2@f110.dev").Get(m, f.Proxy.URL("test"))
					})

					s.It("should not proxy to the backend", func(m *framework.Matcher) bool {
						return m.Equal(http.StatusUnauthorized, m.LastResponse().StatusCode)
					})
				})
			})
		})

		s.Context("access to the backend which via connector", func(s *framework.Scenario) {
			var api1, api2 *framework.MockServer
			s.BeforeAll(func(m *framework.Matcher) bool {
				api1 = f.Proxy.MockServer(m, "api1")
				api2 = f.Proxy.MockServer(m, "api2")
				f.Proxy.Backend(&configv2.Backend{
					Name: "test",
					HTTP: []*configv2.HTTPBackend{
						{Path: "/api1", Agent: true},
						{Path: "/api2", Agent: true},
					},
					Permissions: []*configv2.Permission{
						{Name: "all", Locations: []configv2.Location{{Any: "/"}}},
					},
				})
				f.Proxy.Connector("test/api1", api1)
				f.Proxy.Connector("test/api2", api2)
				f.Proxy.Role(&configv2.Role{Name: "test", Bindings: []*configv2.Binding{{Backend: "test", Permission: "all"}}})
				f.Proxy.User(&database.User{Id: "test@f110.dev", Roles: []string{"test"}})
				return true
			})
			s.AfterAll(func(m *framework.Matcher) bool { return f.Proxy.ClearConf() })

			s.Context("by authenticated user", func(s *framework.Scenario) {
				s.Subject(func(m *framework.Matcher) bool {
					return f.Agents.Authorized("test@f110.dev").Get(m, f.Proxy.URL("test", "/api1"))
				})

				s.It("should proxy to backend", func(m *framework.Matcher) bool {
					m.Equal(http.StatusOK, m.LastResponse().StatusCode)
					m.Len(api1.Requests(), 1)
					return m.Len(api2.Requests(), 0)
				})
			})
		})
	})
}
