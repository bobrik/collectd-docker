#!/bin/sh

set -e

if [ -n "${MESOS_HOST}" ]; then
  if [ -z "${COLLECTD_HOST}" ]; then
    export COLLECTD_HOST="${MESOS_HOST}"
  fi
fi

export GRAPHITE_PORT=${GRAPHITE_PORT:-2003}
export GRAPHITE_PREFIX=${GRAPHITE_PREFIX:-collectd.}
export COLLECTD_INTERVAL=${COLLECTD_INTERVAL:-10}

# Adding a user if needed to be able to communicate with docker
GROUP=nobody
if [ -e /var/run/docker.sock ]; then
  GROUP=$(ls -l /var/run/docker.sock | awk '{ print $4 }')
fi
useradd -g "${GROUP}" collectd-docker-collector

exec reefer -t /etc/collectd/collectd.conf.tpl:/tmp/collectd.conf \
  collectd -f -C /tmp/collectd.conf "$@" > /dev/null
