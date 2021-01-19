# CloudQuery powered by Osquery

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
  -  example [extension_config.json.sample](extension/extension_config.json.sample) is given below.

```json
{
  "aws": {
    "accounts": [
      {
        "id": "12712753535",
        "credentialFile": "/cloudquery/config/credentials",
        "profileName": "default"
      }
    ]
  },
  "gcp": {
    "accounts": [
      {
        "keyFile": "/cloudquery/config/your-serviceAccount.json"
      }
    ]
  },
  "azure": {
    "accounts": [
      {
        "subscriptionId": "3636-3322-dddd-sss-2343444",
        "tenantId": "2377-456-123-266-128635",
        "authFile": "/cloudquery/config/my.auth"
      }
    ]
  }
}

```

- If using aws, copy your aws credentials:
  - `cp ~/.aws/credentials ~/config`
  - Edit credentialFile field  under aws section inside ~/config/extension_config.json and set to /cloudquery/config/credentials
  - Edit id field under aws section inside ~/config/extension_config.json and set to your account id
  - Edit profileName  field under aws section inside ~/config/extension_config.json and set to your  profile name
  - Guide to create AWS credentials: https://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html

- If using Google Cloud, copy your json key file your-serviceAccount.json (cloud be any name) for your service account to `~/config`
  - `cp ~/your-serviceAccount.json ~/config`
  - Edit keyFile field under gcp section inside ~/config/extension_config.json and set to /cloudquery/config/your-serviceAccount.json
  - Guide to create GCP credentials: https://cloud.google.com/iam/docs/creating-managing-service-account-keys

- If using Azure, copy the my.auth (cloud be any name) file for you account to `~/config`
  - `cp ~/my.auth ~/config`
  - Edit authFile  field under azure section inside ~/config/extension_config.json and set to /cloudquery/config/my.auth
  - Edit subscriptionId and tenantId fields under azure section inside ~/config/extension_config.json and set to actual values
  - Guide to create Azure credentials: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest


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
        "keyFile": "/cloudquery/config/your-serviceAccount.json"
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

### Clodquery with osqueryd Docker Container

#### Repeat Configuration under `"Create cloud configurations directory"`

And identify list of scheduled queries and their intervals and place them in `osqyery.conf` inside ~/config on the host. Example osquery.conf is given below.

```json
{
  "schedule": {
    "GCP_COMP_NET": {
      "query": "SELECT * FROM  gcp_compute_network;",
      "interval": 120
    },
    "AWS_S3_BUCK": {
      "query": "SELECT * FROM aws_s3_bucket;",
      "interval": 120
    },
    "AZURE_COMPUTE_VM": {
      "query": "SELECT * FROM azure_compute_vm;",
      "interval": 120
    }
  }
}
```


Once all all the required files under config, run the following commands.

`mkdir ~/query-results` on your host

`sudo docker run -d --rm -v ~/config:/cloudquery/config -v ~/query-results:/var/log/osquery --name cloudquery uptycsdev/cloudconnector:t7 osqueryd`

Now query results can be seen in ~/query-results
