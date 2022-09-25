VERSION=0.1.05

USER_GH=eyedeekay
packagename=i2p-traymenu

GO_COMPILER_OPTS = -a -tags "netgo" -ldflags '-w -extldflags=-static'
WIN_GO_COMPILER_OPTS = -a -tags "netgo windows" -ldflags '-H=windowsgui'

echo:
	@echo "type make version to do release $(VERSION)"

readme:
	grep -v curl README.md | tee README.md.in
	echo "\`\`\`curl -s https://github.com/eyedeekay/i2p-traymenu/releases/download/v$(VERSION)/install.sh | sh\`\`\`" | tee -a README.md.in
	cp README.md.in README.md

version:
	gothub release -p -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION) -d "version $(VERSION)"

del:
	gothub delete -s $(GITHUB_TOKEN) -u $(USER_GH) -r $(packagename) -t v$(VERSION)

tar:
	tar --exclude .git \
		--exclude .go \
		--exclude bin \
		--exclude examples \
		-cJvf ../$(packagename)_$(VERSION).orig.tar.xz .

all: windows osx linux

windows: fmt
	CC=x86_64-w64-mingw32-gcc-win32 CGO_ENABLED=1 GOOS=windows go build $(WIN_GO_COMPILER_OPTS) -o $(packagename).exe
	#CC=i686-w64-mingw32-gcc-win32 CGO_ENABLED=1 GOOS=windows GOARCG=i386 go build $(WIN_GO_COMPILER_OPTS) -o $(packagename)-32.exe

osx: fmt
	#GOARCH=386 GOOS=darwin go build $(GO_COMPILER_OPTS) -o $(packagename)-darwin-386
	GOOS=darwin go build $(GO_COMPILER_OPTS) -o $(packagename)-darwin

linux: fmt
	GOOS=linux go build $(GO_COMPILER_OPTS) -o $(packagename)

sumwindows=`sha256sum $(packagename).exe`
sumlinux=`sha256sum $(packagename)`
sumdarwin=`sha256sum $(packagename)-darwin`

upload-windows:
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumwindows)" -n "$(packagename).exe" -f "$(packagename).exe"

upload-darwin:
	#gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumdarwin)" -n "$(packagename)-darwin" -f "$(packagename)-darwin"

upload-linux:
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumlinux)" -n "$(packagename)" -f "$(packagename)"

upload: upload-windows upload-darwin upload-linux

release: version upload

fmt:
	gofmt -w -s main.go

curlpipe:
	@echo '#! /usr/bin/env sh' | tee install.sh
	@echo "#!/bin/sh" | tee -a install.sh
	@echo 'case "$(uname -s)" in' | tee -a install.sh
	@echo '' | tee -a install.sh
	@echo '   Darwin)' | tee -a install.sh
	@echo "     if [ -f $(packagename) ]; then" | tee -a install.sh
	@echo "       curl -o $(packagename) https://github.com/eyedeekay/i2p-traymenu/releases/download/v$(VERSION)/i2p-traymenu-darwin" | tee -a install.sh
	@echo "     fi" | tee -a install.sh
	@echo '     ;;' | tee -a install.sh
	@echo '' | tee -a install.sh
	@echo '   Linux)' | tee -a install.sh
	@echo "     if [ -f $(packagename) ]; then" | tee -a install.sh
	@echo "       curl -o $(packagename) https://github.com/eyedeekay/i2p-traymenu/releases/download/v$(VERSION)/i2p-traymenu" | tee -a install.sh
	@echo "     fi" | tee -a install.sh
	@echo '     ;;' | tee -a install.sh
	@echo '' | tee -a install.sh
	@echo '   *)' | tee -a install.sh
	@echo '     echo "This system unsupported by curlpipe install"' | tee -a install.sh
	@echo '     ";;"' | tee -a install.sh
	@echo 'esac' | tee -a install.sh
	@echo "sudo chmod a+x $(packagename)" | tee -a install.sh
	@echo "./$(packagename)" | tee -a install.sh

sumpipe=`sha256sum $(packagename)`

upload-pipe: curlpipe readme
	gothub upload -R -u eyedeekay -r "$(packagename)" -t v$(VERSION) -l "$(sumpipe)" -n "curlpipe to install" -f "install.sh"

