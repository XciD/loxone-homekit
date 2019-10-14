FROM golang:1.12-stretch as builder

ENV GO111MODULE on
WORKDIR /go/src/loxone-homekit

COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

FROM alpine as release

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder /go/src/loxone-homekit/build/loxone-homekit /loxone-homekit

ENTRYPOINT ["/loxone-homekit"]