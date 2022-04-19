build:

	go build -o platform-changelog-api cmd/api/main.go

lint:

	gofmt -l .
	gofmt -s -w .

test:

	go test -p 1 -v ./...

run-migration:

	go build -o platform-changelog-migration internal/migration/main.go
	./platform-changelog-migration

compose:

	podman-compose -f development/compose.yml up

compose-quiet:

	podman-compose -f development/compose.yml up -d

compose-down:

	podman-compose -f development/compose.yml down