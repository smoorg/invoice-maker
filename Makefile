PKG=pkg/
BIN=bin/invoice-maker
 
build:
	# added .. as BIN would be created in src/bin otherwise
	go build -o $(BIN) --pkgdir $(PKG) cmd/v2/main.go
 
run:
	go run cmd/v2/main.go

test:
	go test ./...

clean:
	go clean
	rm ./${BIN}
