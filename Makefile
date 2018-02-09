SOURCES = $(wildcard **/*.go)

.PHONY: test
test:
	go test ./...

terraform-provider-auth0: $(SOURCES)
	go build ./...

.PHONY: build
build: terraform-provider-auth0

.PHONY: install
install: terraform-provider-auth0
	mv terraform-provider-auth0 $(shell dirname $(shell which terraform))
