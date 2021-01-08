all: build
BIN-DIR=bin
EXTENSION-DIR=extension

build: extension
extension: $(shell find . -type f)
	go build -o ${BIN-DIR} ./${EXTENSION-DIR}
clean:
	rm -rf ${BIN-DIR}/*