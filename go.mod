module github.com/kyma-project/rafter

go 1.15

require (
	github.com/asyncapi/converter-go v0.0.0-20190916120412-39eeca5e9df5
	github.com/asyncapi/parser v0.0.0-20191002092055-f7b577d06d20
	github.com/gernest/front v0.0.0-20181129160812-ed80ca338b88
	github.com/go-ini/ini v1.51.0 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/minio/minio-go v6.0.14+incompatible
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/vrischmann/envconfig v1.2.0
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	gopkg.in/ini.v1 v1.48.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20190905181640-827449938966 // indirect
	k8s.io/api v0.17.8
	k8s.io/apimachinery v0.17.8
	k8s.io/client-go v0.17.8
	sigs.k8s.io/controller-runtime v0.5.9
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

replace github.com/smartystreets/goconvey => github.com/m00g3n/goconvey v1.6.5-0.20200622160247-ef17e6397c60
