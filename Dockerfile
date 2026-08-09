FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /recommendation-service \
    ./cmd/api

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /recommendation-service /recommendation-service

EXPOSE 8080

ENTRYPOINT ["/recommendation-service"]