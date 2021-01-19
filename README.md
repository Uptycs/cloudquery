# CloudQuery

CloudQuery is OsQuery extension to fetch cloud telemetry from AWS, GCP, and Azure. It is extensible so that  
one can add support for new tables easily, and configurable so that one can change the table schema as well.

## Getting started

### Build

- Checkout the code
- Install prerequisites
  - `go 1.13`
- Set environment varibale for extension home (it shoud be path-to-repo/cloudquery/extension)  
  `export CLOUDQUERY_EXT_HOME=/home/user/work/code/cloudquery/extension`
- Build extension binary.  
  `make`
- To install at default osquery directory (`/etc/osquery/`), run:  
  `make install`

### Test

#### With osqueryi

- Start osqueryi  
  `osqueryi --nodisable_extensions`
- Note down the socket path  
  `.socket`
- `cp ${CLOUDQUERY_EXT_HOME}/extension_config.json.sample ${CLOUDQUERY_EXT_HOME}/extension_config.json`
- Edit `${CLOUDQUERY_EXT_HOME}/extension_config.json` with your cloud accounts. You can add multiple accounts for each cloud provider
- Start extension  
  `./bin/extension --socket /path/to/socket --home-directory ${CLOUDQUERY_EXT_HOME}`
- Query data  
  `select account_id, region_code,image_id,image_type from aws_ec2_image;`

#### With osquery

- Build and install cloudquery
- Edit (or create if does't exist) file `/etc/osquery/extensions.load` and add the following line:
- `/etc/osquery/cloudquery.ext`
- Add following lines to `/etc/osquery/osquery.flags`  
  `--disable_extensions=false`  
  `--extensions_autoload=/etc/osquery/extensions.load`  
  `--extensions_timeout=3`  
  `--extensions_interval=3`
- Copy extension config file to `/etc/osquery/cloudquery`
  - `sudo cp ${CLOUDQUERY_EXT_HOME}/extension_config.json.sample /etc/osquery/cloudquery/extension_config.json`
- Edit `/etc/osquery/cloudquery/extension_config.json` with your cloud accounts. You can add multiple accounts for each cloud provider
  - `sudo vi /etc/osquery/cloudquery/extension_config.json`
- Restart osquery service.
  - `sudo service osqueryd restart`

### Clodquery with osqueryi Docker Container

#### Create cloud configurations directory

- Create a config directory on host to hold the credentials for your cloud accounts (~/config is an example, but this could be any directory):


  - `mkdir ~/config` on the machine where docker container is started
  - ~/config from the host would be mounted to /cloudquery/config inside container 
- Copy `extension_config.json.sample` to your new config directory on your host:
  - `cp extension/extension_config.json.sample ~/config/extension_config.json`

- If using aws, copy your aws credentials:
  - `cp ~/.aws/credentials ~/config`
  - Edit credentialFile field  under aws section inside ~/config/extension_config.json and set to /cloudquery/config/credentials

- If using Google Cloud, copy your json key file cloudConnectorsTrial-serviceAccount.json (cloud be any name) for your service account to `~/config`
  - `cp ~/cloudConnectorsTrial-serviceAccount.json ~/config`
  - Edit keyFile field under gcp section inside ~/config/extension_config.json and set to /cloudquery/config/cloudConnectorsTrial-serviceAccount.json 

- If using Azure, copy the my.auth (cloud be any name) file for you account to `~/config`
  - `cp ~/my.auth ~/config`
  - Edit authFile  field under azure section inside ~/config/extension_config.json and set to /cloudquery/config/my.auth


- After  editing, your  ~/config/extension_config.json  would be looking like as following

```json

{
  "aws": {
    "accounts": [
      {
        "id": "xxxxxxxxxxxx",
        "credentialFile": "/cloudquery/config/credentials",
        "profileName": "default"
      }
    ]
  },
  "gcp": {
    "accounts": [
      {
        "keyFile": "/cloudquery/config/cloudConnectorsTrial-serviceAccount.json"
      }
    ]
  },
  "azure": {
    "accounts": [
      {
        "subscriptionId": "dfffe-3322-dddd-sss-2343444",
        "tenantId": "3333-dfs-333-sfe-121124",
        "authFile": "/cloudquery/config/my.auth"
      }
    ]
  }
}

```




#### Run container with osqueryi

`sudo docker run -it --rm -v ~/config:/cloudquery/config --name cloudquery uptycsdev/cloudconnector:t7`

Press enter to get osquery prompt

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
