module github.com/decentrio/sorobook-api

go 1.24

toolchain go1.24.11

require (
	github.com/decentrio/xdr-converter v1.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	github.com/rakyll/statik v0.1.7
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142
	google.golang.org/grpc v1.64.1
	google.golang.org/protobuf v1.34.2
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.8
)

replace github.com/decentrio/xdr-converter => ../xdr-converter

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stellar/go-stellar-sdk v0.0.0-20251210134752-6c46f8811c13 // indirect
	github.com/stellar/go-xdr v0.0.0-20231122183749-b53fb00bcac2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
)
