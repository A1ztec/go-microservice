FROM alpine:latest

WORKDIR /app

COPY .env .

COPY . /app

CMD ["./loggerApp"]
