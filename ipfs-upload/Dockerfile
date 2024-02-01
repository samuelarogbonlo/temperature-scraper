
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

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

ENV IPFS_API_URL ""
ENV IPFS_GATEWAY_URL ""
ENV DB_HOST ""
ENV DB_PORT ""
ENV DB_PASS ""
ENV DB_NAME ""

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
