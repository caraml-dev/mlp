module github.com/caraml-dev/mlp

go 1.20

require (
	cloud.google.com/go/storage v1.29.0
	github.com/antihax/optional v1.0.0
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/gavv/httpexpect/v2 v2.15.0
	github.com/getsentry/raven-go v0.2.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator v9.30.0+incompatible
	github.com/go-playground/validator/v10 v10.11.1
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/google/go-cmp v0.5.9
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/vault/api v1.9.0
	github.com/hashicorp/vault/api/auth/gcp v0.4.0
	github.com/heptiolabs/healthcheck v0.0.0-20180807145615-6ff867650f40
	github.com/iancoleman/strcase v0.2.0
	github.com/jinzhu/copier v0.3.5
	github.com/jinzhu/gorm v1.9.11
	github.com/knadh/koanf v1.4.4
	github.com/lib/pq v1.3.0
	github.com/newrelic/go-agent v3.19.2+incompatible
	github.com/opentracing/opentracing-go v1.1.0
	github.com/ory/keto-client-go v0.11.0-alpha.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.1
	github.com/prometheus/client_model v0.2.0
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.1
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	go.uber.org/zap v1.17.0
	golang.org/x/exp v0.0.0-20230713183714-613f0c0eb8a1
	golang.org/x/oauth2 v0.5.0
	golang.org/x/sync v0.1.0
	google.golang.org/api v0.106.0
	gopkg.in/errgo.v2 v2.1.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/client-go v0.26.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go v0.107.0 // indirect
	cloud.google.com/go/compute v1.14.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.8.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v3 v3.0.0 // indirect
	github.com/certifi/gocertifi v0.0.0-20191021191039-0944d244cd40 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.1 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.6.6 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sanity-io/litter v1.5.5 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-lib v2.0.0+incompatible // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.34.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yalp/jsonpath v0.0.0-20180802001716-5cc68e5049a0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/term v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/time v0.0.0-20220210224613-90d013bbcef8 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230131230820-1c016267d619 // indirect
	google.golang.org/grpc v1.52.3 // indirect
	google.golang.org/protobuf v1.29.0 // indirect
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apimachinery v0.26.0 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d // indirect
	moul.io/http2curl/v2 v2.3.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
