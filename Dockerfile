FROM golang:1.18.2 as builder

ENV CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN go build -ldflags '-w -s'


FROM alpine:3.16

ENV PRIVATEBIN_URL=https://privatebin.net

WORKDIR /app

# COPY cmds.yaml cmds.yaml

COPY --from=builder /build/balena-data-extractor .

ENTRYPOINT ["./balena-data-extractor"]
