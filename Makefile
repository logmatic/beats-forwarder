PROJECT_NAME=beats-forwarder
BEAT_DIR=github.com/logmatic
GOPACKAGES=$(shell glide novendor)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PREFIX?=.


.PHONY: init
init:
	glide update --strip-vcs
#	make update
	git init

.PHONY: commit
commit:
	git add README.md CONTRIBUTING.md
	git commit -m "Initial commit"
	git add LICENSE
	git commit -m "Add the LICENSE"
	git add .gitignore .gitattributes
	git commit -m "Add git settings"
	git add .
#	git reset -- .travis.yml
	git commit -m "Add beats-forwarder"
#	git add .travis.yml
#	git commit -m "Add Travis CI"

.PHONY: update-deps update
update-deps:
	glide update --strip-vcs

# Checks project and source code if everything is according to standard
.PHONY: check
check:
	@gofmt -l ${GOFILES_NOVENDOR} | read && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true
	go vet ${GOPACKAGES}
