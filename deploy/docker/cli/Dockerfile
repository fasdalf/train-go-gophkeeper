FROM golang:1.22 as local_app

RUN apt-get -y update && apt install -y protobuf-compiler && \
    go install github.com/go-delve/delve/cmd/dlv@v1.23.0 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 && \
    go install github.com/goware/modvendor@v0.5.0
# ##@@ modvendor is for github.com/marcboeker/go-duckdb \
# remove it with duckdb

WORKDIR $GOPATH/src/app/

CMD /bin/bash