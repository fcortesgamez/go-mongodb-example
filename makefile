NO_VENDOR := `go list ./... | grep -v 'vendor/'`

build: test
	go build cmd/webshopd/webshopd.go

test:
	go test github.com/fcortesgamez/go-mongodb-example/internal/...
	go test github.com/fcortesgamez/go-mongodb-example/cmd/...

sanitize:
	go fmt $(NO_VENDOR)
	go vet $(NO_VENDOR)

clean:
	rm -f webshopd
