FROM debian:sid-slim

RUN mkdir /service && apt-get update && apt-get -y install ca-certificates && rm -rf /var/cache/apt/*

COPY ./bin/linkfetcher /service/linkfetcher

RUN useradd server

USER server

EXPOSE 8080

CMD ["/service/linkfetcher"]
