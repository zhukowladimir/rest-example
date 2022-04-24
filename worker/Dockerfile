FROM golang:1.18

WORKDIR /usr/local/app

COPY . .

RUN go mod init github.com/zhukowladimir/rest-example/worker
RUN go mod tidy
RUN go build -o /worker

CMD ["bash", "-c", "while ! curl -s localhost:15672 > /dev/null; do echo waiting for rabbitmq; sleep 3; done; /worker"]
