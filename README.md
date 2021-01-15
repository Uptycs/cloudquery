# CloudQuery

CloudQuery is OsQuery extension to fetch cloud telemetry from AWS, GCP, and Azure. It is extensible to that 
one can add support for new tables easily, and configurable so that one can change the table schema as well.

## Getting started

### Build
- Checkout the code
- Set environment varibale for extension home (it shoud be path-to-repo/cloudquery/extension) 
`export CLOUDQUERY_EXT_HOME=/home/apatil/work/code/cloudquery/extension`
- Copy extension/extension_config.json.sample as extension/extension_config.json and add configurations for
your cloud accounts. You can add multiple accounts for each cloud provider
- `make`

### Test
#### With osqueryi
- Start osqueryi
`osqueryi  --nodisable_extensions`
- Start extension
./bin/extension --socket /home/xyz/.osquery/shell.em
- Query data
`select account_id, region_code,image_id,image_type from aws_ec2_image;`
#### With osquery
TODO

### Supported tables
#### AWS
- aws_ec2_image
- aws_ec2_instance
- aws_ec2_subnet
- aws_ec2_vpc
- aws_s3_bucket

#### GCP
- gcp_compute_disk
- gcp_compute_instance
- gcp_compute_network
- gcp_storage_bucket

#### Azure
- azure_compute_networkinterface
- azure_compute_vm

### Re-configuring a table
TODO
