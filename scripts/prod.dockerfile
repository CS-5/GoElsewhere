FROM alpine:latest

WORKDIR /root

COPY GoElsewhere .
COPY /web/build ./public

CMD ./GoElsewhere
