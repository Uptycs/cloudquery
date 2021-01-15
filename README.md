# CloudQuery

CloudQuery is OsQuery extension to fetch cloud telemetry from AWS, GCP, and Azure. It is extensible to that 
one can add support for new tables easily, and configurable so that one can change the table schema as well.

## Getting started

### Build
- Checkout the code
- Set environment varibale for extension home (it shoud be path-to-repo/cloudquery/extension) 
`export CLOUDQUERY_EXT_HOME=/home/apatil/work/code/cloudquery/extension`
- Copy extension/extension_config.json.sample as CLOUDQUERY_EXT_HOME/extension_config.json and add configurations for
your cloud accounts. You can add multiple accounts for each cloud provider
- `make`
- To install at default osquery directory (/etc/osquery/), run: `make install`

### Test
#### With osqueryi
- Start osqueryi
`osqueryi  --nodisable_extensions`
- Start extension
./bin/extension --socket /home/xyz/.osquery/shell.em
- Query data
`select account_id, region_code,image_id,image_type from aws_ec2_image;`
#### With osquery
- Build and install cloudquery
- Create a file /etc/osquery/extensions.load and add following line to it:
- `/etc/osquery/cloudquery.ext`
- Add following lines to /etc/osquery/osquery.flags
`--disable_extensions=false
--extensions_autoload=/etc/osquery/extensions.load
--extensions_timeout=3
--extensions_interval=3`
- Copy extension/extension_config.json.sample as /etc/osquery/cloudquery/extension_config.json and add configurations for
your cloud accounts. You can add multiple accounts for each cloud provider
- Restart osquery service. `sudo service osqueryd restart`

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
