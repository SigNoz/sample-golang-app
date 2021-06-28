# Bookstore REST API using Gin and Gorm

Read the [article](https://blog.logrocket.com/how-to-build-a-rest-api-with-golang-using-gin-and-gorm/).


## To run and configure app to send data to SigNoz:
```
SERVICE_NAME=goApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend:4317> go run main.go
```
*<IP of SigNoz backend:4317> should be without http/https scheme. Eg localhost:4317*

This runs the gin application at port `8090`. Try accessing API at `http://localhost:8090/books`

Below are the apis available to play around. The API calls will generate telemetry data which will be sent to SigNoz which can be viewed at `<IP of SigNoz backend:3000`
```
GET    /books                    
GET    /books/:id               
POST   /books                    
PATCH  /books/:id                
DELETE /books/:id                
```