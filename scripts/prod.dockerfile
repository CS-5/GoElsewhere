FROM alpine:latest

WORKDIR /root

COPY goelsewhere .
COPY /web/build ./public

CMD ./goelsewhere
