all: build
BIN-DIR=bin
EXTENSION-DIR=extension

build: extension
extension: $(shell find . -type f)
	mkdir -p ${BIN-DIR}
	go build -o ${BIN-DIR} ./${EXTENSION-DIR}
install:
	mkdir -p /etc/osquery/cloudquery/aws/ec2
	mkdir -p /etc/osquery/cloudquery/aws/s3
	mkdir -p /etc/osquery/cloudquery/gcp/compute
	mkdir -p /etc/osquery/cloudquery/gcp/storage
	mkdir -p /etc/osquery/cloudquery/azure/compute
	cp ${BIN-DIR}/extension /etc/osquery/cloudquery.ext
	cp extension/aws/ec2/table_config.json /etc/osquery/cloudquery/aws/ec2
	cp extension/aws/s3/table_config.json /etc/osquery/cloudquery/aws/s3
	cp extension/gcp/compute/table_config.json /etc/osquery/cloudquery/gcp/compute
	cp extension/gcp/storage/table_config.json /etc/osquery/cloudquery/gcp/storage
	cp extension/azure/compute/table_config.json /etc/osquery/cloudquery/azure/compute
clean:
	rm -rf ${BIN-DIR}/*