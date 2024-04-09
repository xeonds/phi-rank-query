NAME=phi-rank-query
BINDIR=build
VERSION=1.0.0
BUILDTIME=$(shell date -u)
APK_NAME="../build/Phigros_3.6.0.APK"
GOBUILD=go mod tidy && go build -ldflags '-s -w -X "main.version=$(VERSION)" -X "main.buildTime=$(BUILDTIME)"'
FRONTBUILD=cd web && pnpm i && pnpm run build --outDir=../$(BINDIR)/dist --emptyOutDir

.PHONY: init web

all: linux-amd64 windows-amd64 web

web:
	$(FRONTBUILD)

linux-amd64: 
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./$(BINDIR)/$(NAME)-$@

windows-amd64: 
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o ./$(BINDIR)/$(NAME)-$@.exe

run:
	cd $(BINDIR) && ./$(NAME)-linux-amd64

init:
	(go mod tidy) &\
	(cd web && pnpm i)

unpack:
	cd script && ./unpack.sh $(APK_NAME)

deploy: linux-amd64 unpack web
	docker-compose up -d

clean:
	rm -rf $(BINDIR)/$(NAME)-*
