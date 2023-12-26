module github.com/kneadCODE/crazycat/apps/golib

go 1.21

require (
	github.com/99designs/gqlgen v0.17.41
	github.com/getsentry/sentry-go/otel v0.25.0
	github.com/go-chi/chi/v5 v5.0.11
	github.com/stretchr/testify v1.8.4
	github.com/vektah/gqlparser/v2 v2.5.10
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.44.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.21.0
	go.opentelemetry.io/otel/metric v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/sdk/metric v1.21.0
	go.opentelemetry.io/otel/trace v1.21.0
	go.uber.org/automaxprocs v1.5.3
	go.uber.org/zap v1.26.0
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/getsentry/sentry-go v0.25.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sosodev/duration v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.17.0
	golang.org/x/net v0.0.0-20180724234803-3673e40ba225 => golang.org/x/net v0.17.0
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d => golang.org/x/net v0.17.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd => golang.org/x/net v0.17.0
	golang.org/x/net v0.0.0-20190311183353-d8887717615a => golang.org/x/net v0.17.0
	golang.org/x/net v0.0.0-20221002022538-bcab6841153b => golang.org/x/net v0.17.0
	golang.org/x/sys v0.0.0-20180830151530-49385e6e1522 => golang.org/x/sys v0.15.0
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => golang.org/x/sys v0.15.0
	golang.org/x/sys v0.14.0 => golang.org/x/sys v0.15.0
	golang.org/x/text v0.3.0 => golang.org/x/text v0.14.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c => gopkg.in/yaml.v3 v3.0.1
)
