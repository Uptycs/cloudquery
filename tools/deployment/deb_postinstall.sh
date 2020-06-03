#!/bin/bash

systemctl enable cloudqueryd 2>/dev/null || update-rc.d cloudqueryd defaults

if [ -x /usr/bin/uptycs_audit_conf.sh ]
then
  /usr/bin/uptycs_audit_conf.sh
  rm -f /usr/bin/uptycs_audit_conf.sh
fi

rm -rf /var/cloudquery/cloudquery.db || true
cloudqueryd --flagfile=/etc/cloudquery/cloudquery.flags --clean_database || true
service cloudqueryd start
