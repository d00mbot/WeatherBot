FROM golang:1.19-alpine3.17 AS builder

COPY . /subscription-bot
WORKDIR /subscription-bot

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /subscription-bot/.bin/bot .
COPY --from=0 /subscription-bot/config /config

EXPOSE 8080

CMD ["./bot"]
