
prepare:
	dep ensure -v

build: prepare
	go build -o ./bin/auto-staging-tower -v -ldflags "-X github.com/auto-staging/tower/config.commitHash=`git rev-parse HEAD` -X github.com/auto-staging/tower/config.buildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/auto-staging/tower/config.branch=`git rev-parse --abbrev-ref HEAD` -X github.com/auto-staging/tower/config.version=`git describe --abbrev=0 --tags`"

tests:
	go test ./... -v

run:
	go run main.go
