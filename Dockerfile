# golang alpine
FROM golang:1.20.0-alpine as builder

ARG TARGETARCH
ARG TARGETOS

RUN apk update

ENV GO111MODULE on
ENV GOPATH /

RUN mkdir /opt/nuts && cd /opt/nuts
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -o /opt/hashicorp-vault-proxy

# alpine
FROM alpine:3.17.2
RUN apk update \
  && apk add --no-cache \
             tzdata \
             curl
COPY --from=builder /opt/hashicorp-vault-proxy /opt/hashicorp-vault-proxy

HEALTHCHECK --start-period=5s --timeout=5s --interval=5s \
    CMD curl -f http://localhost:8210/health || exit 1

EXPOSE 8210
ENTRYPOINT ["/opt/hashicorp-vault-proxy"]