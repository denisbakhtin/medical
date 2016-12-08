#Basic makefile

default: build

build: clean vet
	@go build -o miobalans-go

doc:
	@godoc -http=:6060 -index

lint:
	@golint ./...

debug_server: 
	@reflex -c reflex.conf
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
	@rm -f ./miobalans-go
