module go.f110.dev/heimdallr

go 1.15

require (
	cloud.google.com/go v0.60.0 // indirect
	cloud.google.com/go/storage v1.10.0
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/coreos/prometheus-operator v0.36.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/google/go-cmp v0.5.1
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-github/v32 v32.1.0
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jarcoal/httpmock v1.0.5
	github.com/jetstack/cert-manager v0.12.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/minio/minio-go/v6 v6.0.44
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.4.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/prometheus/procfs v0.0.10 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	go.etcd.io/etcd/v3 v3.3.0-rc.0.0.20200518175753-732df43cf85b
	go.f110.dev/protoc-ddl v0.0.0-20200823082125-88eabc7ca14c
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/zap v1.14.1
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/tools v0.0.0-20200822203824-307de81be3f4 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/api v0.29.0
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.9
	k8s.io/apiextensions-apiserver v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	sigs.k8s.io/yaml v1.2.0
	software.sslmate.com/src/go-pkcs12 v0.0.0-20190322163127-6e380ad96778
)

replace (
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.1.0
	k8s.io/client-go => k8s.io/client-go v0.17.9
)
