FROM golang:1.22


WORKDIR /build

COPY wordchain/go.mod .
COPY wordchain/go.sum .
RUN go mod download

COPY wordchain/ .

RUN go build -o main .

FROM golang:1.22-alpine

COPY /build/main .

EXPOSE 8081

CMD ["./main"]
