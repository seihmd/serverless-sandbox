.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/app app/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: test clean build
	sls deploy --verbose --stage dev

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

package:
	sls package

test:
	go test ./app/...

install_tools:
	GOBIN=$$PWD/bin go get github.com/golang/mock/mockgen@latest


