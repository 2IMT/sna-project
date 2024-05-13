FROM golang:1.22


WORKDIR /build

COPY wordchain/go.mod .
COPY wordchain/go.sum .
RUN go mod download

COPY wordchain/ .

ARG VERSION=dev

RUN go build -o main -ldflags=-X=main.version=${VERSION} main.go 

FROM golang:1.22-alpine

COPY /build/main . --load 

EXPOSE 8081

CMD ["./main"]
