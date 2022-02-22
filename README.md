# Bookstore REST API using Gin and Gorm

Read this [article](https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/) on understand how to build a sample golang app using Gin and Gorm.


## To run and configure app to send data to SigNoz:
```
SERVICE_NAME=goApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend>:4317 go run main.go
```
*<IP of SigNoz backend:4317> should be without http/https scheme. Eg localhost:4317*

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
