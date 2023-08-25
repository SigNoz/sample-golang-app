# Bookstore REST API using Gin and Gorm

Read this [article](https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/) on understand how to build a sample golang app using Gin and Gorm.

## Note:
If you are using the `without-instrumentation` branch to instrument the app following the [blog post present here](https://signoz.io/opentelemetry/go/), after you are done with making the necessary changes in the `main.go` file, run:

```
go mod tidy
```
to download the required packages and populate the go.sum file. 


## To run and configure app to send data to SigNoz:

For SigNoz Cloud:

```bash
SERVICE_NAME=goApp INSECURE_MODE=false OTEL_EXPORTER_OTLP_HEADERS=signoz-access-token=<SIGNOZ-INGESTION-TOKEN> OTEL_EXPORTER_OTLP_ENDPOINT=ingest.{region}.signoz.cloud:443 go run main.go
```

- Update `<SIGNOZ-INGESTION-TOKEN>` with the ingestion token provided by SigNoz
- Update `ingest.{region}.signoz.cloud:443` with the ingestion endpoint of your region. Refer to the table below for the same.

| Region | Endpoint                   |
| ------ | -------------------------- |
| US     | ingest.us.signoz.cloud:443 |
| IN     | ingest.in.signoz.cloud:443 |
| EU     | ingest.eu.signoz.cloud:443 |

For SigNoz OSS:

```
SERVICE_NAME=goApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend>:4317 go run main.go
```

- `<IP of SigNoz backend:4317>` should be without http/https scheme. Eg `localhost:4317`.

---

This runs the gin application at port `8090`. Try accessing API at `http://localhost:8090/books`

Below are the apis available to play around. The API calls will generate telemetry data which will be sent to SigNoz which can be viewed at `<IP of SigNoz backend>:3000`

```
GET    /books                    
GET    /books/:id               
POST   /books                    
PATCH  /books/:id                
DELETE /books/:id                
```

# [Instrumentation Packages](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/README.md#instrumentation-packages)

The [OpenTelemetry registry](https://opentelemetry.io/registry/) is the best place to discover instrumentation packages.
It will include packages outside of this project.

The following instrumentation packages are provided for popular Go packages and use-cases.

| Instrumentation Package | Metrics | Traces |
| :---------------------: | :-----: | :----: |
| [github.com/astaxie/beego](./github.com/astaxie/beego/otelbeego) | ✓ | ✓ |
| [github.com/aws/aws-sdk-go-v2](./github.com/aws/aws-sdk-go-v2/otelaws)|  | ✓ |
| [github.com/bradfitz/gomemcache](./github.com/bradfitz/gomemcache/memcache/otelmemcache) |  | ✓ |
| [github.com/emicklei/go-restful](./github.com/emicklei/go-restful/otelrestful) |  | ✓ |
| [github.com/gin-gonic/gin](./github.com/gin-gonic/gin/otelgin) |  | ✓ |
| [github.com/go-kit/kit](./github.com/go-kit/kit/otelkit) |  | ✓ |
| [github.com/gocql/gocql](./github.com/gocql/gocql/otelgocql) | ✓ | ✓ |
| [github.com/gorilla/mux](./github.com/gorilla/mux/otelmux) |  | ✓ |
| [github.com/labstack/echo](./github.com/labstack/echo/otelecho) |  | ✓ |
| [github.com/Shopify/sarama](./github.com/Shopify/sarama/otelsarama) |  | ✓ |
| [go.mongodb.org/mongo-driver](./go.mongodb.org/mongo-driver/mongo/otelmongo) |  | ✓ |
| [google.golang.org/grpc](./google.golang.org/grpc/otelgrpc) |  | ✓ |
| [gopkg.in/macaron.v1](./gopkg.in/macaron.v1/otelmacaron) |  | ✓ |
| [host](./host) | ✓ |  |
| [net/http](./net/http/otelhttp) | ✓ | ✓ |
| [net/http/httptrace](./net/http/httptrace/otelhttptrace) |  | ✓ |
| [runtime](./runtime) | ✓ |  |

### Follow Opentelemetry docs for examples on latest otel releases: [Opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp/otlptrace)
