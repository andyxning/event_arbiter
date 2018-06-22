VERSION := 1.0.0
COMMIT_HASH := ` git rev-parse --short HEAD `

dep:
	which godep || go get github.com/tools/godep 

test: build
	go test -race -v ./...

vet:
	go list ./... | grep -v "./vendor*" | xargs go vet

fmt:
	find . -type f -name "*.go" | grep -v "./vendor*" | xargs gofmt -s -w

build: dep fmt vet
	godep go install -v ./...
	godep go build -v -ldflags "-X github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf.version=$(VERSION) -X github.com/gigalixir/eventarbiter/cmd/eventarbiter/conf.commitHash=$(COMMIT_HASH)" -o eventarbiter github.com/gigalixir/eventarbiter/cmd/eventarbiter

clean:
	rm eventarbiter

.PHONY: fmt test dep build clean vet
