FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor \
    -o main ./cmd/app

FROM alpine:3.19 AS runner

RUN apk --no-cache add ca-certificates tzdata curl && \
    adduser -D appuser

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docker-entrypoint.sh .

# Copy migrations if they exist
COPY --from=builder /app/db/migration/ /app/db/migration/

RUN chmod +x /app/main /app/docker-entrypoint.sh && \
    if [ -f "/app/db/migration/goose" ]; then chmod +x /app/db/migration/goose; fi && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl --silent --fail http://localhost:8080/health || exit 1

ENTRYPOINT ["/bin/sh", "/app/docker-entrypoint.sh"]