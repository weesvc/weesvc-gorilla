# multi-stage build for GoLang tiny deployment
FROM golang:1.20-alpine AS build
ARG build_version=0.0.1-SNAPSHOT
ARG build_revision=unknown
RUN apk --no-cache add build-base git
COPY . $GOPATH/src/weesvc-gorilla/
WORKDIR $GOPATH/src/weesvc-gorilla
RUN make BUILD_VERSION=$build_version BUILD_REVISION=$build_revision setup build

# final build artifact
FROM alpine
RUN apk update && apk add ca-certificates
COPY --from=build /go/src/weesvc-gorilla/bin/weesvc /app/
COPY config-docker.yaml /etc/weesvc/config.yaml
CMD ["/app/weesvc", "version"]
