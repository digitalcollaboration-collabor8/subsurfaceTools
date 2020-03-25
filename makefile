GOCMD=go 
GOINSTALL=$(GOCMD) install 
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"


build:
			$(info **** Building Subsurface Collabor8 clients ****)
			@$(GOINSTALL) ./...
release:	linux windows
linux: 
			$(info **** Building Subsurface Collabor8 linux 64 bit client ****)
			GOOS=linux GOARCH=amd64 go install ${LDFLAGS} ./...
windows:
			$(info **** Building Subsurface Collabor8 windows 64 bit client ****)
			GOOS=windows GOARCH=amd64 go install ${LDFLAGS} ./...