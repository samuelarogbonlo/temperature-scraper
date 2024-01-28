#  GO_VERSION=1.21.3
FROM golang:1.21.3-bullseye AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x


RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server .

FROM alpine:latest AS final


RUN --mount=type=cache,target=/var/cache/apk \
    rm /etc/apt/sources.list.d/cuda*.list && \
    apk --update add \
        ca-certificates \
        tzdata \
        && update-ca-certificates \
        && apt-get clean \
        && rm -rf /var/lib/apt/lists/*

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

ENV KAFKA_SERVER_ADDRESS ""

# Expose the port that the application listens on.
EXPOSE 8080
EXPOSE 9092

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
