export GOPATH         ?= $(PWD)/go
export TARGET_GOOS    ?= linux
export TARGET_GOARCH  ?= amd64

GB         = $(GOPATH)/bin/gb
EXECUTABLE = bin/golang-test-task-$(TARGET_GOOS)-$(TARGET_GOARCH)

build: $(GB)
	GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) $(GB) build

test: $(GB)
	gb test #-test.timeout 1s

clean:
	rm -rfv go bin pkg

$(GB):
	go get -v github.com/constabulary/gb/...

export DOCKER_NAMESPACE ?= souz9/golang-test-task
export DOCKER_NO_CACHE  ?= false

docker-build:
	docker build --no-cache=${DOCKER_NO_CACHE} -t ${DOCKER_NAMESPACE} .
