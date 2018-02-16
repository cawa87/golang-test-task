FROM golang:1.9.4
RUN mkdir -p /app
COPY main.go /app
WORKDIR /app
RUN go get golang.org/x/net/html && go build -o main main.go
CMD ["/app/main"]