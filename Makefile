export GOPATH ?= $(PWD)/go

export TARGET_GOOS   ?= linux
export TARGET_GOARCH ?= amd64
export TARGET_GOARM  ?= 7

GB         = $(GOPATH)/bin/gb
EXECUTABLE = bin/test-$(TARGET_GOOS)-$(TARGET_GOARCH)

build: $(GB)
	GOOS=$(TARGET_GOOS) GOARCH=$(TARGET_GOARCH) GOARM=$(TARGET_GOARM) $(GB) build

test: $(GB)
	gb test

clean:
	rm -rfv go bin pkg

$(GB):
	go get -v github.com/constabulary/gb/...
