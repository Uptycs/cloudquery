all: build
BIN-DIR=bin
EXTENSION-DIR=extension
INSTALL-DIR?=/etc/cloudquery

build: extension

extension: $(shell find . -type f)
	mkdir -p ${BIN-DIR}
	go build -o ${BIN-DIR} ./${EXTENSION-DIR}

install:
	mkdir -p ${INSTALL-DIR}/aws/ec2
	mkdir -p ${INSTALL-DIR}/aws/s3
	mkdir -p ${INSTALL-DIR}/gcp/compute
	mkdir -p ${INSTALL-DIR}/gcp/storage
	mkdir -p ${INSTALL-DIR}/gcp/iam
	mkdir -p ${INSTALL-DIR}/gcp/sql
	mkdir -p ${INSTALL-DIR}/gcp/dns
	mkdir -p ${INSTALL-DIR}/gcp/file
	mkdir -p ${INSTALL-DIR}/gcp/container
	mkdir -p ${INSTALL-DIR}/gcp/function
	mkdir -p ${INSTALL-DIR}/azure/compute
	mkdir -p ${INSTALL-DIR}/config
	cp ${BIN-DIR}/extension ${INSTALL-DIR}/cloudquery.ext
	cp extension/aws/ec2/table_config.json ${INSTALL-DIR}/aws/ec2
	cp extension/aws/s3/table_config.json ${INSTALL-DIR}/aws/s3
	cp extension/gcp/compute/table_config.json ${INSTALL-DIR}/gcp/compute
	cp extension/gcp/storage/table_config.json ${INSTALL-DIR}/gcp/storage
	cp extension/gcp/iam/table_config.json ${INSTALL-DIR}/gcp/iam
	cp extension/gcp/sql/table_config.json ${INSTALL-DIR}/gcp/sql
	cp extension/gcp/dns/table_config.json ${INSTALL-DIR}/gcp/dns
	cp extension/gcp/file/table_config.json ${INSTALL-DIR}/gcp/file
	cp extension/gcp/container/table_config.json ${INSTALL-DIR}/gcp/container
	cp extension/gcp/function/table_config.json ${INSTALL-DIR}/gcp/function
	cp extension/azure/compute/table_config.json ${INSTALL-DIR}/azure/compute

test:
	@set -x; \
	cd ${EXTENSION-DIR}; \
	go test -v ./...
	@set -x; \
	cd utilities; \
	go test -v ./...

clean:
	rm -rf ${BIN-DIR}/*

.PHONY: all
