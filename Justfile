update:
	go get -u ./...
	go mod tidy
	cd mailing/templating && cargo update
	cd web/engine && cargo update
	cd compiler && cargo update

build-rust:
	cd mailing/templating && cargo build --release
	cd web/engine && cargo build --release
	cd compiler && cargo build --release

build: build-rust
	go build -o jinya-cms-goji app.go

compile-themes:
	for theme in themes/*; do \
		if [ -f "$$theme/Cargo.toml" ]; then \
			cd "$$theme" && cargo build --target wasm32-wasip2 --release; \
		fi \
	done
