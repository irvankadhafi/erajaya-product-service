FROM golang:1.18-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o erajaya

EXPOSE 4000

ENTRYPOINT ["/app/erajaya","migrate"]
ENTRYPOINT ["/app/erajaya","server"]
