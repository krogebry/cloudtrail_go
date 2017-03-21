
dc-build: build
	docker-compose build

deps:
	go get gopkg.in/mgo.v2

build:
	go build -o bin/queue queue/*.go
	go build -o bin/scrape scrape/*.go
	chmod 755 bin/*
