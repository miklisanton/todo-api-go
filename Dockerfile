FROM golang:latest as builder

WORKDIR /build

ADD go.mod .env config.yaml .

COPY . .

RUN go build -o main ./cmd/server

FROM golang:latest

WORKDIR /build

COPY --from=builder /build/main /build/main

COPY --from=builder /build/.env /build/.env

COPY --from=builder /build/config.yaml /build/config.yaml

COPY --from=builder /build/internal/db/migrations /build/internal/db/migrations

CMD ["./main"]

