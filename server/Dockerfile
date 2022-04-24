
FROM golang:1.18

WORKDIR /usr/local/app/store

COPY store/ .

WORKDIR /usr/local/app/server

COPY server/ .

RUN go mod init github.com/zhukowladimir/rest-example/worker
RUN go mod tidy
RUN go build -o /server

CMD ["bash", "-c", "while ! curl -s localhost:15672 > /dev/null; do echo waiting for rabbitmq; sleep 3; done; /server"]
