## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

ORAS ?= $(LOCALBIN)/oras

ORAS_VERSION ?= v1.0.0

.PHONY: oras
oras: $(ORAS) ## Download oras locally if necessary.
$(ORAS): $(LOCALBIN)
	test -s $(LOCALBIN)/oras || GOBIN=$(LOCALBIN) go install oras.land/oras/cmd/oras@$(ORAS_VERSION)

.PHONY: build-amd64
build-mdbook-gh-project-amd64:  ## Build binary for amd64.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/mdbook-gh-project cmd/mdbook-gh-project/main.go

BIN_REGISTRY ?= ghcr.io
BIN_MDBOOK_GH_PROJECT_URL ?= ghcr.io/githedgehog/mdbook-gh-project

.PHONY: bin-registry-login
bin-registry-login: oras
	$(ORAS) login -u "$(USERNAME)" -p "$(PASSWORD)" $(BIN_REGISTRY)

.PHONY: bin-mdbook-gh-project-push
bin-mdbook-gh-project-push: build-mdbook-gh-project-amd64 oras
	cd bin && $(ORAS) push $(BIN_MDBOOK_GH_PROJECT_URL):latest mdbook-gh-project