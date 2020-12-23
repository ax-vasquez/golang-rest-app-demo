include .env

all: clean build

clean:
	@echo
	@echo "[INFO] Clean"
	rm -rf build

build:
	go build -o build/codingtest cmd/main.go

run: build
	build/codingtest

lint:
	@echo
	@echo "[INFO] Get golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$GOPATH/bin v1.24.0
	@echo "[INFO] Running golint"
	$$GOPATH/bin/golangci-lint run --timeout=10m ./...

vet:
	@echo
	@echo "[INFO] Running go vet"
	go vet ./...

secure:
	@echo
	@echo "[INFO] Running gosec"
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.4.0
	$$GOPATH/bin/gosec `ls -d */  | grep -v '.cache' | xargs printf -- '%s... '`

test:
	@echo
	@echo "[INFO] Running tests"
	go test -v -coverprofile=coverage.out ./...

show-coverage:
	@echo
	@echo "[INFO] Opening browser"
	go tool cover -html=coverage.out

race:
	@echo
	@echo "[INFO] Running race condition detection tests"
	go test -race -short ./...
