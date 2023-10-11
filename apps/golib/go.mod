module github.com/kneadCODE/crazycat/apps/golib

go 1.20

require (
	github.com/getsentry/sentry-go v0.19.0
	github.com/getsentry/sentry-go/otel v0.19.0
	github.com/newrelic/go-agent/v3 v3.20.4
	github.com/stretchr/testify v1.8.2
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/sdk v1.14.0
	go.opentelemetry.io/otel/trace v1.14.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.7.0
	golang.org/x/net v0.0.0-20180724234803-3673e40ba225 => golang.org/x/net v0.8.0
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d => golang.org/x/net v0.8.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd => golang.org/x/net v0.8.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a => golang.org/x/net v0.8.0
	golang.org/x/net v0.0.0-20221002022538-bcab6841153b => golang.org/x/net v0.8.0
	golang.org/x/sys v0.0.0-20180830151530-49385e6e1522 => golang.org/x/sys v0.6.0
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => golang.org/x/sys v0.6.0
	golang.org/x/text v0.3.0 => golang.org/x/text v0.8.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c => gopkg.in/yaml.v3 v3.0.1
)
