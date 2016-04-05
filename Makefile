#Basic makefile

default: build

build: vet
	@go generate ./... && go build -o miobalans-go

doc:
	@godoc -http=:6060 -index

lint:
	@golint ./...

debug: clean
	@reflex -c reflex.conf

run: build
	./miobalans-go

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@rm -f ./miobalans-go
	@rm -f ./system/*.rice-box.go
