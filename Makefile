BINARY_NAME=invoice-maker
SRC=src/
PKG=pkg/
# added .. as BIN would be created in src/bin otherwise
BIN=../bin/${BINARY_NAME}
 
build:
	go build -o $(BIN) --pkgdir $(PKG) -C $(SRC) main.go
 
run:
	go build -o $(BIN) --pkdir $(PKG) -C $(SRC) main.go
	./bin/${BINARY_NAME}
 
test:
	go test $(SRC) ./...

clean:
	go clean
	rm ${BINARY_NAME}
