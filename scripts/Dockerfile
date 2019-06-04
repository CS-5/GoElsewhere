# Build Go App
FROM golang:alpine
RUN apk add --no-cache git gcc g++
ENV CGO_ENABLED=1
ENV GOOS=linux
COPY *.go go.mod go.sum ./
RUN go build -a -installsuffix cgo -ldflags "-w -s" -o goelsewhere .

# Build Node App
FROM node:alpine
COPY web/src ./src
COPY web/public ./public
COPY web/package.json web/package-lock.json ./
RUN npm install --production
RUN npm run build

# Build Docker Image
FROM alpine:latest
WORKDIR /app
COPY --from=0 goelsewhere .
COPY --from=1 build ./public

CMD ./goelsewhere
