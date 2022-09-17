.ONESHELL:

build-ui:
	cd webapp/ui && \
	yarn install && \
	yarn build --dest ../site

build-backend: test
	env GOOS=linux GOARCH=arm go build -o piHatDraw .
	tar -czvf piHatDraw-arm.tar.gz piHatDraw
	env GOOS=linux GOARCH=arm64 go build -o piHatDraw .
	tar -czvf piHatDraw-arm64.tar.gz piHatDraw

build: build-ui build-backend

test:
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	ginkgo -r .

.PHONY: build \
        test \
        build-backend \
        build-ui
