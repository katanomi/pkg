module github.com/katanomi/pkg

go 1.23.0

toolchain go1.24.0

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/blendle/zapdriver v1.3.1
	github.com/caarlos0/env/v6 v6.6.2
	github.com/cloudevents/sdk-go/v2 v2.15.2
	github.com/davecgh/go-spew v1.1.1
	github.com/emicklei/go-restful-openapi/v2 v2.3.0
	github.com/emicklei/go-restful/v3 v3.12.2
	github.com/go-logr/logr v1.4.2
	github.com/go-logr/zapr v1.3.0
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.6.0
	github.com/jarcoal/httpmock v1.0.8
	github.com/onsi/ginkgo/v2 v2.19.0
	github.com/onsi/gomega v1.33.1
	github.com/opencontainers/image-spec v1.1.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.12.0
	golang.org/x/time v0.6.0
	k8s.io/api v0.29.15
	k8s.io/apimachinery v0.29.15
	k8s.io/apiserver v0.29.15
	k8s.io/client-go v0.29.15
	knative.dev/pkg v0.0.0-20231023150739-56bfe0dd9626
	sigs.k8s.io/controller-runtime v0.17.6
	sigs.k8s.io/yaml v1.4.0
	yunion.io/x/log v0.0.0-20201210064738-43181789dc74 // indirect
	yunion.io/x/pkg v0.0.0-20210218105412-13a69f60034c
)

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/exporters/zipkin v1.2.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/trace v1.32.0
)

require (
	github.com/alessio/shellescape v1.4.1
	github.com/docker/distribution v2.8.2+incompatible
	github.com/joho/godotenv v1.4.0
	github.com/joshdk/go-junit v1.0.0
	github.com/k1LoW/duration v1.2.0
	github.com/minio/minio-go/v7 v7.0.47
	github.com/mitchellh/mapstructure v1.5.0
	github.com/moby/patternmatcher v0.5.0
	github.com/open-policy-agent/opa v0.61.0
	github.com/opencontainers/go-digest v1.0.0
	github.com/pkg/errors v0.9.1
	// try to keep this the same version as controller-runtime
	github.com/prometheus/client_golang v1.18.0
	github.com/regclient/regclient v0.7.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.6
	github.com/tektoncd/pipeline v0.56.0
	github.com/tidwall/gjson v1.14.4
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	golang.org/x/net v0.37.0
	k8s.io/apiextensions-apiserver v0.29.15
	k8s.io/cli-runtime v0.29.15
	k8s.io/klog/v2 v2.130.1
	oras.land/oras-go/v2 v2.3.1
)

require (
	contrib.go.opencensus.io/exporter/ocagent v0.7.1-0.20200907061046-05415f1de66d // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230305170008-8188dc5388df // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.8.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/cel-go v0.18.1 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-containerregistry v0.17.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.1.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.32.0 // indirect
	go.starlark.net v0.0.0-20230525235612-a134d8f9ddca // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/automaxprocs v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/api v0.156.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.29.15 // indirect
	k8s.io/kube-openapi v0.0.0-20250304201544-e5f78fe3ede9 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/kustomize/api v0.13.5-0.20230601165947-6ce0bf390ce3 // indirect
	sigs.k8s.io/kustomize/kyaml v0.14.3-0.20230601165947-6ce0bf390ce3 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/openzipkin/zipkin-go v0.4.0 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8
)

replace (
	go.uber.org/zap => github.com/katanomi/zap v1.27.0-katanomi // indirect
	knative.dev/pkg => github.com/katanomi/knative-pkg v0.0.0-20240409030957-1aad1aa6fd94

	// branch: alaudx/release-0.17
	sigs.k8s.io/controller-runtime => github.com/katanomi/controller-runtime v0.17.6-1
)
