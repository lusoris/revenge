module github.com/lusoris/revenge

go 1.26.0

require (
	// Validation
	github.com/go-playground/validator/v10 v10.30.1
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/confmap v1.0.0
	github.com/knadh/koanf/providers/env v1.1.0
	github.com/knadh/koanf/providers/file v1.2.1

	// Configuration
	github.com/knadh/koanf/v2 v2.3.2
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/alexedwards/argon2id v1.0.0
	github.com/asticode/go-astiav v0.40.0
	github.com/aws/aws-sdk-go-v2 v1.41.1
	github.com/aws/aws-sdk-go-v2/config v1.32.7
	github.com/aws/aws-sdk-go-v2/credentials v1.19.7
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.22.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.96.0
	github.com/casbin/casbin/v2 v2.135.0
	github.com/coreos/go-oidc/v3 v3.17.0
	github.com/davidbyttow/govips/v2 v2.16.1-0.20250707035900-51ebed754616
	github.com/fergusstrange/embedded-postgres v1.33.0
	github.com/go-faster/errors v0.7.1
	github.com/go-faster/jx v1.2.0
	github.com/go-webauthn/webauthn v0.15.0
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/google/uuid v1.6.0
	github.com/govalues/decimal v0.1.36
	github.com/hashicorp/go-hclog v1.6.3
	github.com/hashicorp/raft v1.7.3
	github.com/hashicorp/raft-boltdb/v2 v2.3.1
	github.com/imroc/req/v3 v3.57.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/lmittmann/tint v1.1.3
	github.com/maypok86/otter/v2 v2.3.0
	github.com/ogen-go/ogen v1.18.0
	github.com/pquerna/otp v1.5.0
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/rueidis v1.0.71
	github.com/redis/rueidis/mock v1.0.71
	github.com/redis/rueidis/rueidisotel v1.0.71
	github.com/riverqueue/river v0.30.2
	github.com/riverqueue/river/riverdriver/riverpgxv5 v0.30.2
	github.com/riverqueue/river/rivertype v0.30.2
	github.com/sony/gobreaker v1.0.0
	github.com/stretchr/testify v1.11.1
	github.com/testcontainers/testcontainers-go v0.40.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.40.0
	github.com/typesense/typesense-go/v2 v2.0.0
	go.opentelemetry.io/otel v1.40.0
	go.opentelemetry.io/otel/metric v1.40.0
	go.opentelemetry.io/otel/trace v1.40.0
	go.uber.org/fx v1.24.0
	go.uber.org/mock v0.6.0
	golang.org/x/oauth2 v0.35.0
	golang.org/x/sync v0.19.0
	golang.org/x/time v0.14.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/asticode/go-astikit v0.42.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.6 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/casbin/govaluate v1.3.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/docker/docker v28.5.1+incompatible // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-webauthn/x v0.1.26 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/go-tpm v0.9.6 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7 // indirect
	github.com/hashicorp/go-immutable-radix v1.0.0 // indirect
	github.com/hashicorp/go-metrics v0.5.4 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.2 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/hashicorp/raft-boltdb v0.0.0-20251103221153-05f9dd7a5148 // indirect
	github.com/icholy/digest v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/lib/pq v1.11.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.1.0 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.57.1 // indirect
	github.com/refraction-networking/utls v1.8.1 // indirect
	github.com/riverqueue/river/riverdriver v0.30.2 // indirect
	github.com/riverqueue/river/rivershared v0.30.2 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/shirou/gopsutil/v4 v4.25.6 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	google.golang.org/grpc v1.78.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
