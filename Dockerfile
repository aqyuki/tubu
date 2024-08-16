FROM golang:1.23.0 AS builder

ARG BOT_VERSION=unknown
ARG BUILD_DATE=unknown
ARG COMMIT_HASH=unknown

WORKDIR /app
RUN --mount=type=bind,target=. go mod download
RUN --mount=type=bind,target=. go mod verify
RUN --mount=type=bind,target=. go build -o /dist/tubu -ldflags="-s -w \
  -X 'github.com/aqyuki/tubu/packages/metadata.Version=${BOT_VERSION}' \
  -X 'github.com/aqyuki/tubu/packages/metadata.GoVersion=$(go version | awk '{print $3}' | sed 's/go//')' \
  -X 'github.com/aqyuki/tubu/packages/metadata.BuildDate=${BUILD_DATE}' \
  -X 'github.com/aqyuki/tubu/packages/metadata.CommitHash=${COMMIT_HASH}'" \
  main.go

FROM gcr.io/distroless/cc-debian12 AS runner

WORKDIR /app
COPY --from=builder --chown=root:root /dist/tubu /app/tubu
STOPSIGNAL SIGINT
ENTRYPOINT ["./tubu"]
