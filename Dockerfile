# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Download dependencies first so this layer is cached when only source changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source tree and compile a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./cmd/api/main.go

# ── Stage 2: run ─────────────────────────────────────────────────────────────
FROM alpine:3.20

# Add CA certificates so the binary can make outbound TLS calls if ever needed.
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]

