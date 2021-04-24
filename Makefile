PROJECTNAME=$(shell basename "$(PWD)")

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

format:
	go fmt .
.PHONY: format

test: format
	@echo "$(PROJECTNAME) tests are running"
	go test .
.PHONY: test

test-coverage: format
	@echo "$(PROJECTNAME) test coverages are running"
	rm -rf .test-coverage
	mkdir .test-coverage
	go test -coverprofile .test-coverage/coverage.out .
	go tool cover -html=.test-coverage/coverage.out -o .test-coverage/cover.html
	open .test-coverage/cover.html
.PHONY: test-coverage

pre-commit-hook:
	printf "make test && make swagger" > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
.PHONY: pre-commit-hook