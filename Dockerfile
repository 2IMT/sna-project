FROM golang:1.22 AS builder

WORKDIR /build

COPY wordchain/go.mod .
COPY wordchain/go.sum .
RUN go mod download

COPY wordchain/ .

ARG VERSION=dev

RUN go build -o main -ldflags=-X=main.version=${VERSION} .

FROM golang:1.22
COPY --from=builder /build/main .

EXPOSE 8081

CMD ["./main"]
