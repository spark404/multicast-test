FROM golang:1.19.0-alpine3.16 as builder
WORKDIR $GOPATH/ns

RUN apk add build-base git cmake make

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/receiver/*.go cmd/receiver/

RUN CGO_ENABLED=0 go build -o /receiver ./cmd/receiver


FROM alpine:3.16 as runner
COPY --from=builder /receiver /usr/bin/receiver

CMD /usr/bin/receiver