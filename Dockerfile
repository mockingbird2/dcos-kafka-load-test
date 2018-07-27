FROM golang:latest

RUN mkdir -p /go/src/dcos-kafka-load-test

ADD . /go/src/dcos-kafka-load-test
WORKDIR /go/src/dcos-kafka-load-test

RUN go build
RUN chmod +x scripts/run.sh
CMD ["/go/src/dcos-kafka-load-test/scripts/run.sh"]
