# Makefile for Go Wechaty
#
# 	GitHb: https://github.com/wechaty/python-wechaty
# 	Author: Huan LI <zixia@zixia.net> git.io/zixia
#

SOURCE_GLOB=$(wildcard bin/*.go src/**/*.go tests/**/*.go examples/*.go)
VERSION=$(shell cat VERSION)

.PHONY: all
all : clean lint

.PHONY: clean
clean:
	rm -fr dist/*
	echo "clean what?"

.PHONY: lint
lint: golint

.PHONY: golint
golint:
	~/go/bin/golint wechaty
	~/go/bin/golint wechaty-puppet
	~/go/bin/golint wechaty-puppet-service

.PHONY: install
install:
	go get -u golang.org/x/lint/golint

.PHONY: gotest
gotest:
	go test `go list ./... | grep -v /vendor/` -v -count=1 -coverpkg=./...

.PHONY: test
test: golint gotest

.PHONY: bot
bot:
	go run examples/ding-dong-bot.go

.PHONY: version
version:
	@newVersion=$$(awk -F. '{print $$1"."$$2"."$$3+1}' < VERSION) \
		&& echo $${newVersion} > VERSION \
		&& echo VERSION := \'$${newVersion}\' > src/version.go \
		&& git add VERSION src/version.py \
		&& git commit -m "$${newVersion}" > /dev/null \
		&& git tag "v$${newVersion}" \
		&& echo "Bumped version to $${newVersion}"
