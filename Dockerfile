FROM ubuntu
MAINTAINER mrsavinov
RUN apt-get update
RUN apt-get install -y software-properties-common python-software-properties
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get -y install golang-go git
RUN mkdir /go
ENV GOPATH=/go
ENV cuncurrentRequests

RUN go get github.com/mrsavinov/golang-test-task
RUN go build github.com/mrsavinov/golang-test-task

CMD ["./golang-test-task", "-bind", "0.0.0.0:8888", "-concurrentRequests", "2000", "-workers", "0"]