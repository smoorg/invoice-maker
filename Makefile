PKG=pkg/
BIN=bin/invoice-maker
 
build:
	# added .. as BIN would be created in src/bin otherwise
	go build -o $(BIN) --pkgdir $(PKG) main.go
 
run:
	go run .

test:
	go test ./...

clean:
	go clean
	rm ./${BIN}
