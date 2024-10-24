#Basic makefile

default: build build_assets

deploy:
	ansible-playbook deploy.yml -K

build: clean vet
	@echo "Building application"
	CGO_ENABLED=0 go build -o miobalans-go

build_assets:
	@echo "Building assets"
	@gulp

watch:
	@air

watch_assets:
	@gulp watch

vet:
	@go vet ./...

clean:
	@echo "Cleaning binary"
	@rm -f ./miobalans-go
