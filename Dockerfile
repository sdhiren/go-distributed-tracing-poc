FROM golang:alpine

WORKDIR /app
COPY . /app

RUN go build -o ./api1/main ./api1
RUN go build -o ./api2/main ./api2
RUN go build -o ./api3/main ./api3
RUN go build -o ./api4/main ./api4

EXPOSE 8080 8081 8082 8083

RUN chmod 755 ./start.sh

CMD ["./api1/main"]
# CMD ["./start.sh"]