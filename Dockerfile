FROM alpine:3.4

RUN apk add --no-cache ca-certificates

WORKDIR /
COPY ./service /service

EXPOSE 8080

CMD ["./service"]