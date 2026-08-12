update:
	go get -u ./...
	go mod tidy
	cd libgoji && cargo update

lint-rust:
	cd libgoji && cargo clippy --release

lint-go:
    go vet .

lint-frontend:
    bun run designer:lint
    bun run installer:lint

lint: lint-rust lint-go lint-frontend

build-rust:
	cd libgoji && cargo build --release

build-go:
    go build -o jinya-cms-goji app.go

build: build-rust build-go

start: build-rust
    go run app.go

compile-themes:
	for theme in themes/*; do \
		if [ -f "$$theme/Cargo.toml" ]; then \
			cd "$$theme" && cargo build --target wasm32-wasip2 --release; \
		fi \
	done

prepare-documentation:
	cd docs && bun install

build-documentation: prepare-documentation
    bunx @redocly/cli bundle openapi/v1/openapi.yaml > docs/apidocs.yml
    rm -rf docs/api
    cd docs && bun docusaurus gen-api-docs all

serve-documentation:
    cd docs && bun start

deploy-documentation: build-documentation
    cd docs && bun run build
    mv docs/build public
