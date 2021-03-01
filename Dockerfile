FROM golang:1.16 as builder

ENV GO111MODULE on
WORKDIR /go/src/loxone-homekit

COPY Makefile go.mod go.sum ./
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
