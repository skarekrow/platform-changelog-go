build:

	go build -o platform-changelog-api cmd/api/main.go

lint:

	gofmt -l .
	gofmt -s -w .

test:

	go test -p 1 -v ./...

run-migration:

	go build -o platform-changelog-migration cmd/migration/main.go
	./platform-changelog-migration