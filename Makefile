APP_NAME=maild
VPATH=.:src:$(GOPATH)/src:

update-deps: deps
	@echo Updating dependencies

build: deps
	go build -o $(APP_NAME) src/main.go

run: deps
	go run src/main.go

clean:
	rm $(APP_NAME)

deps: github.com/lib/pq

github.com/lib/pq:
	go get github.com/lib/pq
