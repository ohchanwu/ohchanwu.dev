# Build stage
FROM golang:1.25.7-alpine3.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static

USER nonroot:nonroot
EXPOSE 8080 8443
ENTRYPOINT ["/app/server"]
