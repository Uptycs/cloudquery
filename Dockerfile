# Copyright (c) 2020-present, The cloudquery authors
#
# This source code is licensed as defined by the LICENSE file found in the
# root directory of this source tree.
#
# SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)

FROM ubuntu:20.04

ARG OSQUERY_VERSION=4.6.0
ARG CLOUDQUERY_VERSION

LABEL \
  name="cloudquery" \
  description="cloudquery powered by Osquery" \
  version="${CLOUDQUERY_VERSION}" \
  url="https://github.com/Uptycs/cloudquery"

ADD https://pkg.osquery.io/deb/osquery_${OSQUERY_VERSION}-1.linux_amd64.deb /tmp/osquery.deb
ADD cloudquery /usr/bin/cloudquery.ext

RUN set -ex; \
    DEBIAN_FRONTEND=noninteractive apt-get update -y && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends curl jq -y && \
    dpkg -i /tmp/osquery.deb && \
    /etc/init.d/osqueryd stop && \
    rm -rf /var/osquery/* /var/log/osquery/* /var/lib/apt/lists/* /var/cache/apt/* /tmp/* && \
    groupadd -g 1000 cloudquery && \
    useradd -m -g cloudquery -u 1000 -d /opt/cloudquery -s /bin/bash cloudquery && \
    mkdir /opt/cloudquery/etc /opt/cloudquery/var && \
    echo "/usr/bin/cloudquery.ext" > /opt/cloudquery/etc/autoload.exts && \
    chmod 700 /usr/bin/cloudquery.ext && \
    chown cloudquery:cloudquery /usr/bin/osquery? /usr/bin/cloudquery.ext /opt/cloudquery/autoload.exts /opt/cloudquery/etc /opt/cloudquery/var

USER cloudquery

ENV CLOUDQUERY_EXT_HOME /opt/cloudquery

WORKDIR /opt/cloudquery

COPY osquery.flags /opt/cloudquery/etc/
COPY osquery.conf  /opt/cloudquery/etc/

COPY extension/aws/ec2/table_config.json        /opt/cloudquery/aws/ec2/
COPY extension/aws/s3/table_config.json         /opt/cloudquery/aws/s3/
COPY extension/gcp/compute/table_config.json    /opt/cloudquery/gcp/compute/
COPY extension/gcp/storage/table_config.json    /opt/cloudquery/gcp/storage/
COPY extension/azure/compute/table_config.json  /opt/cloudquery/azure/compute/

ENTRYPOINT ["/usr/bin/osqueryd", \
    "--flagfile=/opt/cloudquery/etc/osquery.flags", \
    "--config_path=/opt/cloudquery/etc/osquery.conf", \
    "--ephemeral", \
    "--disable_logging", \
    "--database_path=/opt/cloudquery/osquery.db", \
    "--extensions_socket=/opt/cloudquery/var/osquery.em", \
    "--extensions_autoload=/opt/cloudquery/etc/autoload.exts"]
