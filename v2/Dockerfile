FROM alpine:3.17.2 as certs
RUN apk update && apk add ca-certificates

FROM golang:1.20.2-alpine3.17 as builder
ARG COMMAND="nct"
ARG GIT_VER
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/nct-linux-amd64 -ldflags "-s -w -X main.command=${COMMAND} -X main.version=${GIT_VER}"

FROM scratch
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /usr/src/app/build/nct-linux-amd64 /nct
ENTRYPOINT ["/nct"]