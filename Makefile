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

run-api:

	./platform-changelog-api

run-db:

	podman run --rm -it -p 5432:5432 -e POSTGRES_PASSWORD=crc -e POSTGRES_USER=crc -e POSTGRES_DB=gumbaroo --name postgres postgres:12.4

compose:

	podman-compose -f development/compose.yml up

compose-quiet:

	podman-compose -f development/compose.yml up -d

compose-down:

	podman-compose -f development/compose.yml down
