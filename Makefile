#Basic makefile

default: build build_assets

deploy: build build_assets
	./deploy.sh

build: clean vet
	@echo "Building application"
	CGO_ENABLED=0 go build -o medical-go cmd/main.go

build_assets:
	@echo "Building assets"
	@gulp

watch:
	@air

watch_assets:
	@gulp watch

vet:
	@go vet ./...

lint:
	@golangci-lint run

clean:
	@echo "Cleaning binary"
	@rm -f ./medical-go
