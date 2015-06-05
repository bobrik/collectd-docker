#!/bin/sh

set -e

if [ -n "${MESOS_HOST}" ]; then
    if [ -z "${COLLECTD_HOST}" ]; then
        export COLLECTD_HOST="${MESOS_HOST}"
    fi
fi

if [ ! -e "/.initialized" ]; then
    touch "/.initialized"
    envtpl /etc/collectd/collectd.conf.tpl
fi

exec gosu nobody collectd -f > /dev/null
