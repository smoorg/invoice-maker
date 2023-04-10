SRC=src/
PKG=pkg/
BIN=bin/invoice-maker
 
build:
	# added .. as BIN would be created in src/bin otherwise
	go build -o ../$(BIN) --pkgdir $(PKG) -C $(SRC) main.go
 
run:
	cd $(SRC) && go run .

test:
	go test -C $(SRC) ./...

clean:
	cd ${SRC} && go clean
	rm -f ./${BIN}
