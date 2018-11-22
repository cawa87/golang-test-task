build:
	go build -o bin/scrapper ./scrapper

test:
	go test -count=1 ./scrapper/...
