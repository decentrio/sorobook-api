# Sorobook API

This repository offers high-performance APIs with a gRPC server seamlessly integrated with the Sorobook backend, powered by PostgreSQL for efficient data querying. It serves various types of Soroban data in a friendly JSON format.

## Features

We support APIs for querying events, transactions, and ledgers.

### 1. Events

- Stellar asset contract events
  - From or to an address
  - Triggered by a contract
- Event by ID
- Event by contract
- Event at ledger
- Count of events by contract

### 2. Ledger

- Ledger by hash
- Ledger by sequence

### 3. Transaction

- Transaction at ledger
- Transaction by address
- Transaction by contract
- Transaction by contract and address
- Transaction by hash

For details, visit the testing endpoint: <https://sorobook-api.decentrio.ventures/public/#/>

## Prerequisites

1. Install `swagger-combine`:

   ```bash
   npm install --save swagger-combine
   ```

2. Install `statik`:

   ```bash
   go install github.com/rakyll/statik@v0.1.7
   ```

## Deployment

Follow the instructions below:

1. Set up environment variables:

   ```bash
   export READONLY_URL=<postgresql_url>
   ```

2. Build the binary:

   ```bash
   go build main.go
   ```

3. Run the gRPC server and host the public OpenAPI gateway:

   ```bash
   ./main
   ```

You can then query your own Sorobook through [localhost:8080](http://localhost:8080) or use the OpenAPI UI at [localhost:8080/public](http://localhost:8080/public).
