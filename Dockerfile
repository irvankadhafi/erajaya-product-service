FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o erajaya

EXPOSE 4000

CMD ./erajaya
