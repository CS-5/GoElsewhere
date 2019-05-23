FROM golang:1.11.9-alpine3.9

RUN apk add --no-cache git gcc g++
ENV GO111MODULE=on
WORKDIR /app
COPY main.go go.mod go.sum /app/
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o go-elsewhere . 


FROM node

WORKDIR /web
COPY web/src /web/src
COPY web/public /web/public
COPY ["web/package.json", "web/package-lock.json", "./"]
RUN npm install --production && npm run build


FROM alpine:latest

WORKDIR /root/
COPY --from=0 /app/go-elsewhere .
COPY --from=1 /web/build ./public
CMD ./go-elsewhere