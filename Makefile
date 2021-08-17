GOSRC = $(shell find . -type f -name '*.go')

VERSION=v0.0.1

build: haul

haul: $(GOSRC)
	CGO_ENABLED=0 GOOS=linux go build -o haul cmd/haul/haul.go

mac:
	CGO_ENABLED=0 GOOS=darwin go build -o haul cmd/haul/haul.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o haul.exe cmd/haul/haul.go

clean:
	rm -rf haul

test:
	go test -v -timeout 60s -race ./...
