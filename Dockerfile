FROM golang:1.8

ADD . /go/src/github.com/mihis/golang-test-task
WORKDIR /go/src/github.com/mihis/golang-test-task
RUN make

RUN mkdir /etc/golang-test-task
ARG config="config.yaml"
COPY $config /etc/golang-test-task

EXPOSE 8085

ENTRYPOINT ["/go/bin/golang-test-task", "--config", "/etc/golang-test-task/config.yaml"]