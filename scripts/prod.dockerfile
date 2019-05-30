FROM alpine:latest

WORKDIR /root

COPY go-elsewhere .
COPY /web/build ./public

CMD ./go-elsewhere