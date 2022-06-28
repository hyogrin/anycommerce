default: build

GORELEASER := $(shell command -v goreleaser 2> /dev/null)

SWAGGER_DOCS_PORT ?= 9000

.PHONY: build test install release

build:
	env CGO_ENABLED=0 go build -ldflags="-X github.com/anycommerce/user/pkg.Version=$(shell git describe --tags)" .

test:
	@go test ./... -coverprofile=/tmp/go-cover.tmp

release:
ifndef GORELEASER
	$(error "goreleaser not found (`go get -u -v github.com/goreleaser/goreleaser` to fix)")
endif
	$(GORELEASER) --rm-dist

.PHONY: swagger swagger-ui
swagger:
	@swag fmt
	@swag init --parseDependency --parseInternal --parseDepth 2

swagger-ui:
	@echo "API docs preview will be running at http://localhost:$(SWAGGER_DOCS_PORT)"
	@docker run --rm -v $(PWD)/docs:/docs \
		-e 'SWAGGER_JSON=/docs/swagger.yaml' \
		-p $(SWAGGER_DOCS_PORT):8080 \
		swaggerapi/swagger-ui