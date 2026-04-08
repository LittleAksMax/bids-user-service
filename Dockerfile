# Build stage
FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETARCH

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
COPY migrations/ migrations/
COPY main.go main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -installsuffix cgo -o user-service ./

FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/user-service .
COPY --from=builder /app/migrations ./migrations

ENV MODE=production
ENV PORT=8080

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./user-service"]
