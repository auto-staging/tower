
prepare:
	dep ensure -v

build: prepare
	go build -o ./bin/auto-staging-tower -v

tests:
	go test ./... -v

run:
	go run main.go
