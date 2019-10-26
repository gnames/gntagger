GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get
FLAG_MODULE = GO111MODULE=on
VERSION=`git describe --tags`
VER=`git describe --tags --abbrev=0`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.Build=${DATE} \
                  -X main.Version=${VERSION}"


all: install

test: deps install
	$(FLAG_MODULE) go test ./...

deps:
	$(FLAG_MODULE) $(GOGET) github.com/spf13/cobra/cobra@v0.0.5; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/ginkgo/ginkgo@v1.10.2; \
	$(FLAG_MODULE) $(GOGET) github.com/onsi/gomega@v1.7.0; \

build:
	cd gntagger; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) ${LDFLAGS};

release:
	cd gntagger; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/gntagger-${VER}-linux.tar.gz gntagger; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=darwin GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/gntagger-${VER}-mac.tar.gz gntagger; \
	$(GOCLEAN); \
	GO111MODULE=on GOOS=windows GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	zip -9 /tmp/gntagger-${VER}-win-64.zip gntagger.exe; \
	$(GOCLEAN);

install:
	cd gntagger; \
	GO111MODULE=on $(GOINSTALL) ${LDFLAGS};
