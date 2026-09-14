# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# Variables
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
ARG GO_VERSION={{GO_VERSION}}

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# Development stage
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
FROM golang:${GO_VERSION}-alpine AS development
WORKDIR /src/api
RUN go install github.com/cespare/reflex@v0.3.2
CMD ["reflex", "-c", "reflex.conf"]

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# Build API stage
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
FROM golang:${GO_VERSION}-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG CGO_ENABLED=0
ARG GOFLAGS="-buildvcs=false"

ENV CGO_ENABLED=${CGO_ENABLED}
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV GOFLAGS=${GOFLAGS}

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /bin/server ./cmd/server

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
# API Production stage
# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
FROM alpine:3.23 AS production

RUN apk add --no-cache \
    ca-certificates \
    && addgroup -S -g 1001 app \
    && adduser -u 1001 -S app -G app

COPY --from=build /bin/server /usr/local/bin/server
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD ["wget", "-q", "--spider", "http://localhost:8080/health"]

ENTRYPOINT ["/usr/local/bin/server"]
