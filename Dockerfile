FROM golang:1.22

WORKDIR /build

COPY wordchain/go.mod .
COPY wordchain/go.sum .
RUN go mod download

COPY wordchain/ .



RUN go build -o main -ldflags=-X=main.version=${VERSION} main.go 

FROM debian:buster-slim

FROM golang:1.22-alpine

COPY --from=builder /build/main /go/bin/main --load

ENV PATH="/go/bin:${PATH}"

EXPOSE 8081

CMD ["main"]
