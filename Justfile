update:
	go get -u ./...
	go mod tidy
	cd mailing/templating && cargo update
	cd web/engine && cargo update
	cd compiler && cargo update

build:
	cd mailing/templating && cargo build --release
	cd web/engine && cargo build --release
	cd compiler && cargo build --release
	go build -o jinya-cms-goji app.go
