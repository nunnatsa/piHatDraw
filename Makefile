.ONESHELL:

build-ui:
	cd webapp/ui
	yarn install
	yarn build --dest ../site

build-backend: test
	env GOOS=linux GOARCH=arm go build -o piHatDraw-arm .
	env GOOS=linux GOARCH=arm64 go build -o piHatDraw-arm64 .

build: build-ui build-backend

test:
	go test ./...

.PHONY: build \
        test \
        build-backend \
        build-ui
