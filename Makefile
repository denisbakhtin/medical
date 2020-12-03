#Basic makefile

default: build

deploy:
	ansible-playbook deploy.yml -K

build: clean vet
	@echo "Building assets"
	@gulp
	@echo "Building application"
	CGO_ENABLED=0 go build -o miobalans-go

doc:
	@godoc -http=:6060 -index

lint:
	@golint ./...

debug_server: 
	@watcher
debug_assets:
	@gulp watch

#run 'make -j2 debug' to launch both servers in parallel
debug: clean debug_server debug_assets 

run: build
	./medical

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@echo "Cleaning binary"
	@rm -f ./miobalans-go

stop: 
	@echo "Stopping medical service"
	@sudo systemctl stop medical

start:
	@echo "Starting medical service"
	@sudo systemctl start medical

pull:
	@echo "Pulling origin"
	@git pull origin master

pull_restart: stop pull clean build start