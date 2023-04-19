GOLDFLAGS = "-w -s -extldflags '-z now' -X github.com/xujunjie-cover/vtap-cni/versions.BUILDDATE=$(DATE)"

.PHONY: build-go
build-go:
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildmode=pie -o $(CURDIR)/dist/images/vtap-cni -ldflags $(GOLDFLAGS) -v ./cmd/vtap-cni
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildmode=pie -o $(CURDIR)/dist/images/vtap-cni-daemon -ldflags $(GOLDFLAGS) -v ./cmd/vtap-cni-daemon
