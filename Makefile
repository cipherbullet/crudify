build:
	@go build -o bin/crudify

run: build
	@./bin/crudify