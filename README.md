[![Build](https://github.com/Uptycs/cloudquery/workflows/Build/badge.svg?branch=master)](https://github.com/Uptycs/cloudquery/actions?query=workflow%3ABuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/Uptycs/cloudquery)](https://goreportcard.com/report/github.com/Uptycs/cloudquery)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)
---

# cloudquery powered by Osquery

cloudquery is Osquery extension to fetch cloud telemetry from AWS, GCP, and Azure. It is extensible so that
one can add support for new tables easily, and configurable so that one can change the table schema as well.

## Contents

- [Working with extension](#build-and-test-extension)
  * [Build](#build)
  * [Test](#test)
    + [Test with osqueryi](#with-osqueryi)
    + [Test with osqueryd](#with-osqueryd)
- [Working with docker](#test-with-docker)
  * [Setup](#setup-credentials)
  * [Test with osqueryi](#run-osqueryi-from-cloudquery-container)
  * [Test with osqueryd](#run-osqueryd-from-cloudquery-container)
- [Supported tables](#supported-tables)

---

## Build and test extension

### Build

- Install prerequisites
  - make
  - [go](https://golang.org/doc/install#install)
- Set environment variable for extension home (it shoud be path-to-repo/cloudquery/extension)
  ```sh
  export CLOUDQUERY_EXT_HOME=/home/user/work/cloudquery/extension
  ```
- Build the extension:
  ```sh
  make
  ````

### Test

#### With osqueryi

- Copy extension configuration sample file:
  ```sh
  cp ${CLOUDQUERY_EXT_HOME}/extension_config.json.sample ${CLOUDQUERY_EXT_HOME}/config/extension_config.json
  ```
- Edit `${CLOUDQUERY_EXT_HOME}/config/extension_config.json` with your cloud accounts. You can add multiple accounts for each cloud provider. Change logging path and other parameters to suit your needs. Make sure log path is writable.
- Start osqueryi
  ```sh
  osqueryi --nodisable_extensions --extension ${CLOUDQUERY_EXT_HOME}/../cloudquery
  ```
- Query data:
  ```sql
  SELECT account_id, region_code, image_id, image_type FROM aws_ec2_image;
  ```

#### With osqueryd

- Build and install cloudquery:
  ```sh
  make build
  sudo make install
  ```
- Edit (or create) `/etc/osquery/extensions.load` file and append the following line: `/usr/local/bin/cloudquery.ext`
- Edit `/opt/cloudquery/config/extension_config.json` with your cloud accounts. You can add multiple accounts for each cloud provider. Change logging path and other parameters to suit your needs.
- Add following flags to `/etc/osquery/osquery.flags` (your flag file path could be different)
```
--extensions_autoload=/etc/osquery/extensions.load
--disable_extensions=false
```
- Restart osquery service:
  ```sh
  sudo service osqueryd restart
  ```

---

## Test with docker

### Setup credentials

> Setup credentials before proceeding to testing with `osqueryi` or `osqueryd`

- Create a config directory on the host to hold the credentials for your cloud accounts (~/config is an example, but this could be any directory).

- Make a copy of [extension_config.json.sample](extension/extension_config.json.sample) as `extension_config.json` in a directory called `config` (can be anywhere on your machine)
- Copy cloud credentials to the `config` directory
  - For AWS: `$HOME/.aws/credentials`
  - For GCP: `your-serviceAccount.json` or any JSON file that contains GCP credentials
  - For Azure: `my.auth` or any file that holds Azure credentials

- If using AWS cloud, update the following fields in `aws` section in `config/extension_config.json` file:
  - `credentialFile` should be set to `/opt/cloudquery/etc/config/credentials`
  - `id` should match AWS account ID
  - `profileName` should be same as the profile in your `.aws/credentials` file
  - Guide to create AWS credentials: https://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html

- If using Google cloud, update `keyFile` in `gcp` section in `extension_config.json` file. It should be changed to `/opt/cloudquery/etc/config/your-serviceAccount.json` where `your-serviceAccount.json` is the JSON key file that contains GCP credentials
  - Guide to create GCP credentials: https://cloud.google.com/iam/docs/creating-managing-service-account-keys

- If using Azure, update the following fields in `azure` section in `extension_config.json` file:
  - `authFile` should be set to `/opt/cloudquery/etc/config/my.auth`. `my.auth` should be the name of the file that contains your Azure credentials.
  - `subscriptionId` and `tenantId` fields should be changed to values from your Azure account
  - Guide to create Azure credentials: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest

### Run osqueryi inside cloudquery container

```sh
docker run --rm -it --name cloudquery \
  -v <absolute path to host config directory>:/opt/cloudquery/etc/config \
  uptycs/cloudquery:latest \
  osqueryi --extension /usr/local/bin/cloudquery.ext
```

### Run osqueryd from cloudquery container

Following files and directories can be mounted from the host:
- `/opt/cloudquery/logs`              - Directory that contains the logs
- `/opt/cloudquery/etc/osquery.flags` - Osquery flags file
- `/opt/cloudquery/etc/osquery.conf`  - Osquery configuration JSON file
- `/opt/cloudquery/etc/config`        - Directory that contains Cloud provider credentials and cloudquery configuration JSON

Sample Osquery configuration with scheduled queries that can be overwritten via `osquery.conf`:
```json
{
  "schedule": {
    "gcp_compute_network": {
      "query": "SELECT * FROM gcp_compute_network;",
      "interval": 120
    },
    "aws_s3_bucket": {
      "query": "SELECT * FROM aws_s3_bucket;",
      "interval": 120
    },
    "azure_compute_vm": {
      "query": "SELECT * FROM azure_compute_vm;",
      "interval": 120
    }
  }
}
```

```sh
docker run --rm -d --name cloudquery \
  -v <absolute path to host config directory>:/opt/cloudquery/etc/config \
  uptycs/cloudquery:latest
```

---

### Supported tables
- [AWS](extension/aws/tables.md)
- [GCP](extension/gcp/tables.md)
- [Azure](extension/azure/tables.md)
