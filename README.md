# Crawl

# Env

| Name  | Description | Default 
|-------|-------------|---------
| CRAWL_BIND_ADDR         | Bind address        | :9000
| CRAWL_MAX_WORKERS       | Max Scrape Workers  | 100

# Install

``` bash 
mkdir $GOPATH/github.com/l-vitaly/golang-test-task
cd $GOPATH/src/github.com/l-vitaly/golang-test-task
git clone git@github.com:l-vitaly/golang-test-task.git
cd golang-test-task
make dep install
```

# Docker build image

``` bash 
mkdir $GOPATH/github.com/l-vitaly/golang-test-task
cd $GOPATH/src/github.com/l-vitaly/golang-test-task
git clone git@github.com:l-vitaly/golang-test-task.git
cd golang-test-task
make docker c="crawl"
docker run -d -P --name crawl crawl
```
