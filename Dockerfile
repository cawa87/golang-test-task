FROM golang:1.8.3-alpine3.6 as build
RUN apk add --update git make
WORKDIR /workdir
COPY . .
RUN make

FROM alpine:3.6
MAINTAINER Aleksey Nikitin <imsouz9@gmail.com>
RUN apk add --no-cache ca-certificates
COPY --from=build /workdir/bin/golang-test-task-* /bin/golang-test-task
EXPOSE 80
ENTRYPOINT ["/bin/golang-test-task", "-l:80"]
