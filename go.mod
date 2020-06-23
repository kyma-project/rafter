module github.com/kyma-project/rafter

go 1.13

require (
	github.com/asyncapi/converter-go v0.0.0-20190916120412-39eeca5e9df5
	github.com/asyncapi/parser v0.0.0-20191002092055-f7b577d06d20
	github.com/gernest/front v0.0.0-20181129160812-ed80ca338b88
	github.com/go-ini/ini v1.51.0 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/vrischmann/envconfig v1.2.0
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	gopkg.in/ini.v1 v1.48.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20190905181640-827449938966 // indirect
	k8s.io/api v0.0.0-20191003000013-35e20aa79eb8
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20191003000419-f68efa97b39e
	sigs.k8s.io/controller-runtime v0.4.0
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

replace github.com/smartystreets/goconvey => github.com/m00g3n/goconvey v1.6.5-0.20200622160247-ef17e6397c60
