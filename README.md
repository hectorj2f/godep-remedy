# godep-remedy

```
godep-remedy cli

Usage:
  godep-remedy [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Discover go module versions
  help        Help about any command
  update      update a comma-separated list of go modules to update
  version     Version for secret

Flags:
  -h, --help   help for godep-remedy

Use "godep-remedy [command] --help" for more information about a command.
```

`godep-remedy discover` finds the latest available versions for a specific go mod root path.

```
godep-remedy discover .
2023/12/01 17:21:40 Found module::: name: github.com/sirupsen/logrus - from: v1.9.0 - to: v1.9.3
```

`godep-remedy discover --module <gomodulepath>`: list all the available versions for a go module.

```
godep-remedy discover --module github.com/sirupsen/logrus
2023/12/01 17:23:29 Found versions: [v0.1.0 v0.1.1 v0.2.0 v0.3.0 v0.4.0 v0.4.1 v0.5.0 v0.5.1 v0.6.0 v0.6.1 v0.6.2 v0.6.3 v0.6.4 v0.6.5 v0.6.6 v0.7.0 v0.7.1 v0.7.2 v0.7.3 v0.8.0 v0.8.1 v0.8.2 v0.8.3 v0.8.4 v0.8.5 v0.8.6 v0.8.7 v0.9.0 v0.10.0 v0.11.0 v0.11.1 v0.11.2 v0.11.3 v0.11.4 v0.11.5 v1.0.0 v1.0.1 v1.0.3 v1.0.4 v1.0.5 v1.0.6 v1.1.0 v1.1.1 v1.2.0 v1.3.0 v1.4.0 v1.4.1 v1.4.2 v1.5.0 v1.6.0 v1.7.0 v1.7.1 v1.8.0 v1.8.1 v1.8.2 v1.8.3 v1.9.0 v1.9.1 v1.9.2 v1.9.3]
```


`godep-remedy update --module <module@version>,<module@version>...`: update the go mod file to the specific modules and versions listed here.

```
update a comma-separated list of go modules to update

Usage:
  godep-remedy update [gomodule@version] [flags]

Flags:
  -h, --help                      help for update
      --modroot string            go mod root path
      --modules string            comma separated list of go modules with versions
      --remedy                    EXPERIMENTAL: enable auto-solve the go mod conflicts (default true)
      --tidy                      go mod tidy flag
      --validate-module-version   check the module version and error if it founds a newer version (defa
```


```
godep-remedy update --tidy=true  --validate-module-version=false --modules=go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc@v0.46.0,go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.46.0,go.opentelemetry.io/otel/sdk@v1.21.0,go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@v1.21.0
2023/12/01 17:36:43 Updating the modules...
2023/12/01 17:36:43 Remediate with module: go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc and version v0.46.0
2023/12/01 17:36:43 Remediate with module: go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp and version v0.46.0
2023/12/01 17:36:43 Remediate with module: go.opentelemetry.io/otel/sdk and version v1.21.0
2023/12/01 17:36:43 Remediate with module: go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc and version v1.21.0
2023/12/01 17:36:43 Get package: go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
2023/12/01 17:36:43 Running go mod edit require ...
2023/12/01 17:36:43 Get package: go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
2023/12/01 17:36:43 Running go mod edit require ...
2023/12/01 17:36:43 Get package: go.opentelemetry.io/otel/sdk
2023/12/01 17:36:43 Running go mod edit require ...
2023/12/01 17:36:43 Get package: go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
2023/12/01 17:36:43 Running go mod edit require ...
2023/12/01 17:36:43 Running go mod tidy ...
```