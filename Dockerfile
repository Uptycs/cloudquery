FROM ubuntu:20.04 AS build

COPY . /cloudquery/
WORKDIR /cloudquery

RUN set -ex; \
  export DEBIAN_FRONTEND=noninteractive \
  && apt-get update -y \
  && apt-get upgrade -y \
  && apt-get install -y \
  tzdata \
  git \
  golang \
  build-essential \
  && mkdir -p /cloudquery/bin \
  && go build -ldflags="-s -w" -o bin ./extension \
  && make install

FROM ubuntu:20.04

LABEL \
  com.uptycs.description="Uptycs cloudquery container image" \
  com.uptycs.name="cloudquery" \
  com.uptycs.version="1.0" \
  com.uptycs.schema-version="1.0" \
  com.uptycs.url="https://www.uptycs.com" \
  com.uptycs.vendor="Uptycs Inc"
# CHROOT_DIR is the docker directory which points to the host's "/"
# SOFTWARE_DIR will be the location on host where Osquery files will be copied to
ENV \
  CHROOT_DIR="/host" \
  SOFTWARE_DIR="/etc/osquery" \
  CLOUDQUERY_DIR="/cloudquery/extension" \
  OLD_DIR="/cloudquery/extension" \
  CLOUDQUERY_EXT_HOME="/cloudquery/extension" \
  CLOUDQUERY_EXE="/cloudquery/bin/extension"

COPY --from=build /cloudquery/bin/extension ${CLOUDQUERY_EXE}
COPY --from=build ${OLD_DIR}/aws/ec2/table_config.json ${CLOUDQUERY_DIR}/aws/ec2/
COPY --from=build ${OLD_DIR}/aws/s3/table_config.json ${CLOUDQUERY_DIR}/aws/s3/
COPY --from=build ${OLD_DIR}/gcp/compute/table_config.json ${CLOUDQUERY_DIR}/gcp/compute/
COPY --from=build ${OLD_DIR}/gcp/storage/table_config.json ${CLOUDQUERY_DIR}/gcp/storage/
COPY --from=build ${OLD_DIR}/azure/compute/table_config.json ${CLOUDQUERY_DIR}/azure/compute/

    

WORKDIR ${CLOUDQUERY_DIR}
COPY  entrypoint.sh ${CLOUDQUERY_DIR}/
RUN set -ex; \
  DEBIAN_FRONTEND=noninteractive apt-get update -y && \
  DEBIAN_FRONTEND=noninteractive apt-get upgrade -y && \
  apt-get install -y tzdata; \
  apt update;  \
  apt-get install -y wget; \
  wget -O /tmp/osquery.deb  https://pkg.osquery.io/deb/osquery_4.6.0-1.linux_amd64.deb; \
  apt-get update; \
  dpkg -i /tmp/osquery.deb && \
  /etc/init.d/osqueryd stop && \
  mkdir -p /cloudquery/extension/config; \
  chmod 700 ${CLOUDQUERY_EXE}; \
  apt-get purge -y wget; \
  apt-get install -y supervisor; \
  chmod 755 ${CLOUDQUERY_DIR}/entrypoint.sh; \ 
  rm -rf /var/osquery/* /var/log/osquery/* /var/lib/apt/lists/* /var/cache/apt/* /tmp/*;

COPY osquery.flags ${SOFTWARE_DIR}/ 
COPY osqueryd_script.conf /etc/supervisor/conf.d/
ENTRYPOINT ["/cloudquery/extension/entrypoint.sh"]
CMD ["osqueryi", "--nodisable_extensions", "--extension", "/cloudquery/bin/extension"]
