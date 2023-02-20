run:
	go run main.go

build:
	go build

mock:
	mockery --output=./src/mocks/services --dir=./src/services --all --keeptree --case=underscore
	mockery --output=./src/mocks/repositories --dir=./src/repositories --all --keeptree --case=underscore
	mockery --output=./src/mocks/controllers --dir=./src/controllers --all --keeptree --case=underscore

test:
	go test ./... -coverprofile cover.out
