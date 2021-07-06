# ------------------------------------------------------------------------------
# Make build from go sources
# ------------------------------------------------------------------------------
FROM golang:1.16.5-alpine3.14 as build

COPY    go.* /usr/src/syncdbdocs/
WORKDIR /usr/src/syncdbdocs
RUN     go mod download

COPY    . /usr/src/syncdbdocs
RUN     go build

# ------------------------------------------------------------------------------
# Release build will not include source code
# ------------------------------------------------------------------------------
FROM alpine:3.14 as release
COPY --from=build /usr/src/syncdbdocs/syncdbdocs /usr/bin/syncdbdocs

ENTRYPOINT ["/usr/bin/syncdbdocs"]