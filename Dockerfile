# Build stage
FROM golang:alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY service/ service/
COPY repository/ repository/
COPY health/ health/
COPY db/ db/
COPY contracts/ contracts/
COPY config/ config/
COPY cache/ cache/
COPY api/ api/
COPY main.go main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./

FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/user-service .

ENV MODE=production
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./user-service"]
