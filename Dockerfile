# Copyright (c) 2020-present, The cloudquery authors
#
# This source code is licensed as defined by the LICENSE file found in the
# root directory of this source tree.
#
# SPDX-License-Identifier: (Apache-2.0 OR GPL-2.0-only)

FROM ubuntu:20.04

ARG BASEQUERY_VERSION=4.6.0
ARG CLOUDQUERY_VERSION

LABEL \
  name="cloudquery" \
  description="cloudquery powered by Osquery" \
  version="${CLOUDQUERY_VERSION}" \
  url="https://github.com/Uptycs/cloudquery"

ADD https://uptycs-basequery.s3.amazonaws.com/${BASEQUERY_VERSION}/basequery_${BASEQUERY_VERSION}-1.linux_amd64.deb /tmp/basequery.deb
COPY cloudquery /usr/local/bin/cloudquery.ext

RUN set -ex; \
    DEBIAN_FRONTEND=noninteractive apt-get update -y && \
    DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates && \
    dpkg -i /tmp/basequery.deb && \
    /etc/init.d/osqueryd stop && \
    rm -rf /var/osquery/* /var/log/osquery/* /var/lib/apt/lists/* /var/cache/apt/* /tmp/* && \
    groupadd -g 1000 cloudquery && \
    useradd -m -g cloudquery -u 1000 -d /opt/cloudquery -s /bin/bash cloudquery && \
    mkdir /opt/cloudquery/etc /opt/cloudquery/logs /opt/cloudquery/var && \
    echo "/usr/local/bin/cloudquery.ext" > /opt/cloudquery/etc/extensions.load && \
    chmod 700 /usr/local/bin/cloudquery.ext && \
    chown cloudquery:cloudquery /usr/bin/osquery? /usr/local/bin/cloudquery.ext /opt/cloudquery/etc/extensions.load /opt/cloudquery/etc /opt/cloudquery/logs /opt/cloudquery/var

USER cloudquery

ENV CLOUDQUERY_EXT_HOME /opt/cloudquery/etc

WORKDIR /opt/cloudquery

COPY osquery.flags osquery.conf                 /opt/cloudquery/etc/
COPY extension/aws/ec2/table_config.json        /opt/cloudquery/etc/aws/ec2/
COPY extension/aws/iam/table_config.json        /opt/cloudquery/etc/aws/iam/
COPY extension/aws/s3/table_config.json         /opt/cloudquery/etc/aws/s3/
COPY extension/aws/kms/table_config.json        /opt/cloudquery/etc/aws/kms/
COPY extension/aws/cloudtrail/table_config.json /opt/cloudquery/etc/aws/cloudtrail/
COPY extension/aws/cloudfront/table_config.json /opt/cloudquery/etc/aws/cloudfront/
COPY extension/aws/codepipeline/table_config.json /opt/cloudquery/etc/aws/codepipeline/
COPY extension/aws/acm/table_config.json        /opt/cloudquery/etc/aws/acm/
COPY extension/aws/codedeploy/table_config.json /opt/cloudquery/etc/aws/codedeploy/
COPY extension/aws/cloudwatch/table_config.json /opt/cloudquery/etc/aws/cloudwatch/
COPY extension/aws/cloudformation/table_config.json /opt/cloudquery/etc/aws/cloudformation/
COPY extension/aws/apigateway/table_config.json /opt/cloudquery/etc/aws/apigateway/
COPY extension/aws/codecommit/table_config.json /opt/cloudquery/etc/aws/codecommit/
COPY extension/aws/config/table_config.json     /opt/cloudquery/etc/aws/config/
COPY extension/aws/directoryservice/table_config.json  /opt/cloudquery/etc/aws/directoryservice/
COPY extension/aws/guardduty/table_config.json /opt/cloudquery/etc/aws/guardduty/
COPY extension/aws/ecs/table_config.json        /opt/cloudquery/etc/aws/ecs/
COPY extension/aws/ecr/table_config.json     /opt/cloudquery/etc/aws/ecr/
COPY extension/azure/compute/table_config.json  /opt/cloudquery/etc/azure/compute/
COPY extension/aws/organizations/table_config.json /opt/cloudquery/etc/aws/organizations/
COPY extension/aws/workspaces/table_config.json /opt/cloudquery/etc/aws/workspaces/
COPY extension/aws/efs/table_config.json        /opt/cloudquery/etc/aws/efs/
COPY extension/aws/s3_glacier/table_config.json /opt/cloudquery/etc/aws/s3_glacier/
COPY extension/aws/eks/table_config.json /opt/cloudquery/etc/aws/eks/
COPY extension/aws/sns/table_config.json        /opt/cloudquery/etc/aws/sns/
COPY extension/aws/sqs/table_config.json        /opt/cloudquery/etc/aws/sqs/
COPY extension/aws/elb/table_config.json        /opt/cloudquery/etc/aws/elb/
COPY extension/aws/elbv2/table_config.json      /opt/cloudquery/etc/aws/elbv2/
COPY extension/gcp/compute/table_config.json    /opt/cloudquery/etc/gcp/compute/
COPY extension/gcp/dns/table_config.json        /opt/cloudquery/etc/gcp/dns/
COPY extension/gcp/file/table_config.json       /opt/cloudquery/etc/gcp/file/
COPY extension/gcp/iam/table_config.json        /opt/cloudquery/etc/gcp/iam/
COPY extension/gcp/storage/table_config.json    /opt/cloudquery/etc/gcp/storage/
COPY extension/gcp/sql/table_config.json        /opt/cloudquery/etc/gcp/sql/
COPY extension/gcp/container/table_config.json  /opt/cloudquery/etc/gcp/container/
COPY extension/gcp/function/table_config.json   /opt/cloudquery/etc/gcp/function/
COPY extension/gcp/run/table_config.json        /opt/cloudquery/etc/gcp/run/
COPY extension/gcp/cloudlog/table_config.json   /opt/cloudquery/etc/gcp/cloudlog/

CMD ["/usr/bin/osqueryd", \
    "--flagfile=/opt/cloudquery/etc/osquery.flags", \
    "--config_path=/opt/cloudquery/etc/osquery.conf", \
    "--ephemeral", \
    "--logger_path=/opt/cloudquery/logs", \
    "--database_path=/opt/cloudquery/osquery.db", \
    "--extensions_socket=/opt/cloudquery/var/osquery.em", \
    "--extensions_autoload=/opt/cloudquery/etc/extensions.load"]
