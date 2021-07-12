module github.com/katanomi/pkg

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/blendle/zapdriver v1.3.1
	github.com/caarlos0/env/v6 v6.6.2
	github.com/emicklei/go-restful-openapi/v2 v2.3.0
	github.com/emicklei/go-restful/v3 v3.5.1
	github.com/go-logr/logr v0.3.0
	github.com/go-logr/zapr v0.2.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.5
	github.com/onsi/gomega v1.10.2
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.10.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.uber.org/zap v1.17.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.2
	knative.dev/pkg v0.0.0-20210510175900-4564797bf3b7
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/yaml v1.2.0
	yunion.io/x/log v0.0.0-20201210064738-43181789dc74 // indirect
	yunion.io/x/pkg v0.0.0-20210218105412-13a69f60034c
)

replace go.uber.org/zap v1.17.0 => github.com/katanomi/zap v1.18.2 // indirect
