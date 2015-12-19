#!/bin/sh

set -e

if [ -n "${MESOS_HOST}" ]; then
    if [ -z "${COLLECTD_HOST}" ]; then
        export COLLECTD_HOST="${MESOS_HOST}"
    fi
fi

if [ ! -e "/.initialized" ]; then
    envtpl /etc/collectd/collectd.conf.tpl
    touch "/.initialized"
fi

exec gosu nobody collectd -f > /dev/null
