update:
	go get -u ./...
	go mod tidy
	cd mailing/templating && cargo update
	cd web/engine && cargo update
	cd compiler && cargo update

lint-rust:
	cd mailing/templating && cargo clippy --release
	cd web/engine && cargo clippy --release
	cd compiler && cargo clippy --release

lint-go:
    go vet .

lint-frontend:
    bun run designer:lint
    bun run installer:lint

lint: lint-rust lint-go lint-frontend

build-rust:
	cd mailing/templating && cargo build --release
	cd web/engine && cargo build --release
	cd compiler && cargo build --release

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

build-documentation:
    bunx @redocly/cli bundle openapi/v1/openapi.yaml > docs/apidocs.yml
    rm -rf docs/api
    cd docs && bun docusaurus gen-api-docs all

serve-documentation:
    cd docs && bun start

deploy-documentation: build-documentation
    cd docs && bun run build
    mv docs/build public
