module github.com/kyma-project/rafter

go 1.15

require (
	github.com/asyncapi/converter-go v0.3.0
	github.com/asyncapi/parser-go v0.3.0
	github.com/gernest/front v0.0.0-20181129160812-ed80ca338b88
	github.com/go-ini/ini v1.51.0 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/vrischmann/envconfig v1.3.0
	go.uber.org/multierr v1.5.0 // indirect
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	gopkg.in/ini.v1 v1.48.0 // indirect
	k8s.io/api v0.17.11
	k8s.io/apimachinery v0.17.11
	k8s.io/client-go v0.17.11
	sigs.k8s.io/controller-runtime v0.5.10
)

replace (
	github.com/smartystreets/goconvey => github.com/m00g3n/goconvey v1.6.5-0.20200622160247-ef17e6397c60
	go.etcd.io/etcd => go.etcd.io/etcd v3.3.25+incompatible
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.3.0
)
