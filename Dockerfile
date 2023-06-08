ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /build

COPY ./ .
RUN go build -o ./app ./cmd/aeza-promo-instances-watchdog


FROM alpine:3.15

RUN adduser -u 1000 -h /app -D -g "" user  \
    && chown -hR user: /app

WORKDIR /app

COPY --from=builder --chown=user:user /build/app .

USER user

CMD ["./app"]