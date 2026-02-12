# ============================================================
# Multi-stage Dockerfile for Job Scheduler
# ============================================================
# Stage 1 (builder): Build
# Stage 2 (runtime): Minimal runtime image
# ============================================================

# ---- Build Stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Build static binary (CGO disabled, Linux target)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/job-scheduler .

# ---- Runtime Stage ----
# Alpine: minimal footprint (~5MB)
FROM alpine:3.19

# Run as non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy only the built binary (no source in runtime)
COPY --from=builder /app/job-scheduler .

# Switch to non-root user
USER appuser

# Interactive CLI — run with: docker run -it job-scheduler
ENTRYPOINT ["./job-scheduler"]
