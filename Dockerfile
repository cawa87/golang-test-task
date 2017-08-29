FROM ubuntu
MAINTAINER mrsavinov
RUN apt-get update
RUN apt-get install -y software-properties-common python-software-properties
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get -y install golang-go git
RUN mkdir /go
ENV GOPATH=/go

RUN go get github.com/mrsavinov/golang-test-task
RUN go build github.com/mrsavinov/golang-test-task

CMD /golang-test-task -bind $BIND_ADD -concurrentRequests $CONCURRENT_REQUESTS -workers $WORKERS