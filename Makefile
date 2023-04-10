BINARY_NAME=invoice-maker
SRC=src/
PKG=pkg/
# added .. as BIN would be created in src/bin otherwise
BIN=../bin/${BINARY_NAME}
 
build:
	go build -o $(BIN) --pkgdir $(PKG) -C $(SRC) main.go
 
run:
	$(build)
	./bin/${BINARY_NAME}
 
test:
	go test -C $(SRC) ./...

clean:
	go clean
	rm ${BINARY_NAME}
