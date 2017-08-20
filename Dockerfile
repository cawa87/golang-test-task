FROM golang:1.8.3-alpine3.6
MAINTAINER Aleksey Nikitin <imsouz9@gmail.com>

RUN apk add --update git make

WORKDIR /workdir
COPY . .
RUN make && INSTALL_PREFIX=/ make install

RUN apk del git make \
	&& rm -rf /var/cache/apk/* /workdir /go /usr/local/go

WORKDIR /tmp
EXPOSE 80
ENTRYPOINT ["/bin/golang-test-task", "-l:80"]
