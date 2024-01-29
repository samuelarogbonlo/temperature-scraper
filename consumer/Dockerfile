FROM golang:1.21.3-bullseye AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -o /bin/server .

FROM alpine:3.19.1 AS final


RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

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


COPY --from=build /bin/server /bin/

ENV POSTGRES_CONNECTION_STRING ""
ENV KAFKA_SERVER_ADDRESS ""

# Expose the port that the application listens on.
EXPOSE 8081

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
