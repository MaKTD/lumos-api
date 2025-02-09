FROM golang:1.22.12-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o build/main ./cmd/api

FROM golang:1.22.12-alpine3.21 AS runtime
WORKDIR /app
COPY --from=builder /app/build/main ./main

CMD ["./main"]