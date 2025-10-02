# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/main ./cmd/api

# Runtime Stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/out/main .

EXPOSE 3000

CMD ["./main"]