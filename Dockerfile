FROM golang:1.21-alpine AS builder

ARG BARON_CHAIN_VERSION
ENV GOPATH=""
ENV GOMODULE="on"
ENV CGO_ENABLED=0

# Ensure Baron Chain version is specified
RUN test -n "${BARON_CHAIN_VERSION}"

# Install build dependencies
RUN apk add --no-cache make git

# Set working directory
WORKDIR /baron-chain

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY internal/ internal/
COPY testing/ testing/
COPY modules/ modules/
COPY LICENSE .
COPY contrib/devtools/Makefile contrib/devtools/
COPY Makefile .

# Build Baron Chain daemon
RUN --mount=type=cache,target=/root/.cache/go-build \
    make build

# Production image
FROM alpine:3.19

ARG BARON_CHAIN_VERSION
LABEL "org.baron-chain.version" "${BARON_CHAIN_VERSION}"

# Add runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary
COPY --from=builder /baron-chain/build/barond /usr/local/bin/

# Create non-root user
RUN addgroup -g 1000 baron && \
    adduser -D -u 1000 -G baron baron

# Switch to non-root user
USER baron

# Set entrypoint
ENTRYPOINT ["barond"]
