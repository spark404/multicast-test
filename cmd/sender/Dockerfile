FROM golang:1.19.0-alpine3.16 as builder
WORKDIR $GOPATH/ns

RUN apk add build-base git cmake make

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/sender/*.go cmd/sender/

RUN go build CGO_ENABLED=0 -o /sender ./cmd/sender


FROM alpine:3.16 as runner
COPY --from=builder /sender /usr/bin/sender

CMD ["/usr/bin/sender"]