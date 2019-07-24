
prepare:
	dep ensure -v

build: prepare
	GOOS=linux go build -o ./bin/auto-staging-tower -v -ldflags "-X github.com/auto-staging/tower/config.commitHash=`git rev-parse HEAD` -X github.com/auto-staging/tower/config.buildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/auto-staging/tower/config.branch=`git rev-parse --abbrev-ref HEAD` -X github.com/auto-staging/tower/config.version=`git describe --abbrev=0 --tags` -d -s -w" -tags netgo -installsuffix netgo

tests:
	go test ./... -v

run:
	go run main.go
