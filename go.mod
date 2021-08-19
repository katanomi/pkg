module github.com/katanomi/pkg

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/blendle/zapdriver v1.3.1
	github.com/caarlos0/env/v6 v6.6.2
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/emicklei/go-restful-openapi/v2 v2.3.0
	github.com/emicklei/go-restful/v3 v3.5.1
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.2.0
	github.com/go-openapi/spec v0.20.2 // indirect
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.6
	github.com/googleapis/gnostic v0.5.3 // indirect
	github.com/jarcoal/httpmock v1.0.8
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.3
	github.com/opentracing/opentracing-go v1.1.0
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.uber.org/zap v1.18.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	k8s.io/klog/v2 v2.5.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210113233702-8566a335510f // indirect
	knative.dev/pkg v0.0.0-20210730172132-bb4aaf09c430
	sigs.k8s.io/controller-runtime v0.8.3
	sigs.k8s.io/yaml v1.2.0
	yunion.io/x/log v0.0.0-20201210064738-43181789dc74 // indirect
	yunion.io/x/pkg v0.0.0-20210218105412-13a69f60034c
)

replace go.uber.org/zap => github.com/katanomi/zap v1.18.2 // indirect
