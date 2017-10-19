# Pkg path
PKG_PATH=./pkg

# Path to compile
CMD_PATH=./cmd/crawl

# Binary output name
TARGET=crawl

.PHONY: all install uninstall

all: build

dep:
	dep ensure

main.go: ${CMD_PATH}/main.go
	CGO_ENABLED=0 go build -a -installsuffix cgo -o ${TARGET} ${CMD_PATH}

build:	dep main.go

docker:
	docker build -t $$c .

install:
	install ./${TARGET} ${GOPATH}/bin

uninstall:
	rm ${GOPATH}/bin/${TARGET}