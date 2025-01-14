# train-go-gophkeeper
diploma 2

#### running
* docker-compose up -d
* ???
* profit

Migrate:

~~`GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URI" goose up-by-one -dir=internal/server/db/storage/migrations -allow-missing`~~

`docker-compose exec cli go run ./cmd/gophkeeperserver/ migrate`


#### Future roadmap

How to secure gRPC connection with SSL/TLS in Go
https://dev.to/techschoolguru/how-to-secure-grpc-connection-with-ssl-tls-in-go-4ph
https://stackoverflow.com/questions/71876783/make-grpc-call-from-go-client-with-tls-gcp-cloud-function

