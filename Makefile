SHELL=/bin/bash

export GOPATH := $(CURDIR)

.PHONY: all js go

all: js go

go:
	go install -v services/...

js: public/app.js

public/app.js: $(shell find client -name "*.js")
	node_modules/browserify/bin/cmd.js  -t [ babelify --presets [ react ] ] client/main.js > $@
