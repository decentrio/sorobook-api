# Sorobook API

This repository offers a high-performance APIs with a gRPC server seamlessly integrated with sorobook backend which powered by PostgresSQL for efficient data querying. It serves many types of soroban data in friendly json format.

### Features

We support apis for querying events, transactions and ledger 

**1. Events**

- Stellar asset contract events
    - from or to and address
    - trigger by a contract
- Event by ID
- Event by contract
- Event at ledger
- Count of events by contract

**2. Ledger**
- Ledger by hash
- Ledger by seq

**3. Transaction**
- Transaction at ledger
- Transaction by address
- Transaction by contract 
- Transaction by contract and address
- Transaction by hash

For details, visit this testing endpoint: https://sorobook-api.decentrio.ventures/public/#/

### Prerequisite

- Install `swagger-combine`:
```
npm install --save swagger-combine
```
- Install `statik`:
```
go install github.com/rakyll/statik@v0.1.7
```

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
