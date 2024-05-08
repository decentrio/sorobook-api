# Sorobook API

This repository offers a comprehensive solution for building high-performance APIs with a gRPC server seamlessly integrated with a public OpenAPI specification, all powered by PostgreSQL for efficient data querying. Designed with scalability, reliability, and ease of use in mind, this project empowers developers to create robust and efficient systems.

### Deployment
Follow below instructions:
- Setup env:
```
export READONLY_URL=<postgresql_url>
```
- Building binary:
```
go build main.go
```
- Running grpc server and host public OpenAPI gateway:
```
./main
```
Then you can query your own Sorobook through [localhost:8080](localhost:8080) or use OpenAPI UI at [localhost:8080/public](localhost:8080/public)