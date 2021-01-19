#!/bin/bash
set -e
echo "$1" 
echo "daya sankar"





if [ "$1" = osqueryi ]; then
osqueryi --nodisable_extensions --extension /cloudquery/bin/extension
fi

if [ "$1" = osqueryd ]; then
    echo run osqueryd
    if [ -f "/etc/osquery/cloudquery.ext" ]; then
        echo "/etc/osquery/cloudquery.ext exists."
    else 
        echo "/etc/osquery/cloudquery.ext does not exist."
        cp /cloudquery/bin/extension /etc/osquery/cloudquery.ext
    fi
    echo "/etc/osquery/cloudquery.ext" > /etc/osquery/extensions.load  
    #service osqueryd start
    if [ -f "/cloudquery/config/osquery.conf" ]; then
        echo "/cloudquery/config/osquery.conf exists."
        cp /cloudquery/config/osquery.conf /etc/osquery/osquery.conf
    fi
    /usr/bin/supervisord -c /etc/supervisor/conf.d/osqueryd_script.conf 
    #osqueryd_script.conf 
    sleep 150
fi
