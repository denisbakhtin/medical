#Basic makefile

default: build

deploy:
	ansible-playbook deploy.yml -K

build: clean vet
	#@echo "Building assets"
	#@gulp
	@echo "Building application"
	CGO_ENABLED=0 go build -o miobalans-go

watch:
	@gulp watch

vet:
	@go vet ./...

clean:
	@echo "Cleaning binary"
	@rm -f ./miobalans-go
