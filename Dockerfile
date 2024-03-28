FROM alpine:3.14.1 AS base-certs
RUN apk update && apk add --no-cache ca-certificates

FROM scratch AS base
COPY docker-user /etc/passwd
USER 101
COPY --from=base-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

FROM base as nlk
ENTRYPOINT ["/nlk"]
COPY build/nlk /
