FROM  golang:1.25-alpine AS build

WORKDIR /app

# Cache Go module downloads

# Copy source and build
COPY . .

COPY go.mod go.sum ./
RUN go mod tidy
# Build a statically linked binary
# CGO_ENABLED=0 removes C library dependencies
# -ldflags strips debug info for a smaller binary

WORKDIR /app/cmd/api

RUN  go build -o payment-service .

EXPOSE 8080

CMD ["./payment-service"]