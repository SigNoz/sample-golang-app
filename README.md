# Bookstore REST API using Gin and Gorm

Read the [article](https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/).

## How to run:

```
$ go run main.go
```


## To run and configure app to send data to SigNoz:
```
SERVICE_NAME=goApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend:4317> go run main.go
```