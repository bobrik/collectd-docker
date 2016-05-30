# Collect docker container resource usage

This is collectd plugin and docker image to collect resource usage from docker
containers. Resource usage collected from `docker stats` API and sent to
graphite installation. Containers can be added and removed on the fly, no need
to restart collectd.

## Configuration

This plugin treats containers as tasks that run as parts of apps.

### Setting the App name of a Container

* Set the label `collectd_docker_app` directly on the container
* Set `collectd_docker_app_label` on the container that points to which actual
label to use. e.g.`collectd_docker_app_label=app_id` will use `app_id` label on the
container
* Set environment variable `COLLECTD_DOCKER_APP` on the container
* Set `COLLECTD_DOCKER_APP_ENV` on the container that points to which actual
environment variable to use. For example, marathon sets `MARATHON_APP_ID` and
by setting `COLLECTD_DOCKER_APP_ENV` to `MARATHON_APP_ID` you would get the
marathon app id.

These keys can be changed globally by setting `APP_LABEL_KEY` or `APP_ENV_KEY`
when running the collectd container. For example, if you set `APP_ENV_KEY` to
`MARATHON_APP_ID` on the collectd container, then this will use
`MARATHON_APP_ID` on all running containers.

### Setting the Task name of a Container

* Set the label `collectd_docker_task` directly on the container
* Set `collectd_docker_task_label` on the container that points to which actual
label to use. e.g.`collectd_docker_task_label=task_id` will use `task_id` label on the
container
* Set environment variable `COLLECTD_DOCKER_TASK` on the container
* Set `COLLECTD_DOCKER_TASK_ENV` on the container that points to which actual
environment variable to use. For example, mesos sets `MESOS_TASK_ID` and by
setting `COLLECTD_DOCKER_TASK_ENV` to `MESOS_TASK_ID` you would get the mesos
task id.

These keys can be changed globally by setting `TASK_LABEL_KEY` or `TASK_ENV_KEY`
when running the collectd container. For example, if you set `TASK_ENV_KEY` to
`MESOS_TASK_ID` on the collectd container, then this will use `MESOS_TASK_ID` on
all running containers.

### Limitations

* If a container's app name cannot be identified, it will be not monitored. So
if you are not seeing metrics, then it means you must check whether the app
name is configured correctly.
* The string `<app>.<task>` is limited by 63 characters. So it is also useful to
set `COLLECTD_DOCKER_APP_ENV_TRIM_PREFIX` and/or
`COLLECTD_DOCKER_TASK_ENV_TRIM_PREFIX` on the containers.

## Reported metrics

Metric names look line this:

```
collectd.<host>.docker_stats.<app>.<task>.<type>.<metric>
```

Gauges:

* CPU
    * `cpu.user`
    * `cpu.system`
    * `cpu.total`

* Memory overview
    * `memory.limit`
    * `memory.max`
    * `memory.usage`

* Memory breakdown
    * `memory.active_anon`
    * `memory.active_file`
    * `memory.cache`
    * `memory.inactive_anon`
    * `memory.inactive_file`
    * `memory.mapped_file`
    * `memory.pg_fault`
    * `memory.pg_in`
    * `memory.pg_out`
    * `memory.rss`
    * `memory.rss_huge`
    * `memory.unevictable`
    * `memory.writeback`

* Network (bridge mode only)
    * `net.rx_bytes`
    * `net.rx_dropped`
    * `net.rx_errors`
    * `net.rx_packets`
    * `net.tx_bytes`
    * `net.tx_dropped`
    * `net.tx_errors`
    * `net.tx_packets`

Percent:

* CPU
    * `cpu.percent`


## Grafana dashboard

Grafana 2 [dashboard](grafana2.json) is included.

![screenshot](https://github.com/bobrik/collectd-docker/raw/master/screenshot.png)

#### Graphite metrics extracted from the dashboard

* CPU usage per second

```
aliasByNode(scaleToSeconds(nonNegativeDerivative(collectd.$host.docker_stats.$app.$task.gauge.cpu.total), 1), 3, 4, 1)
```

* Memory limit

```
alias(averageSeries(collectd.$host.docker_stats.$app.$task.gauge.memory.limit), 'limit')
```

* Memory usage

```
aliasByNode(collectd.$host.docker_stats.$app.$task.gauge.memory.usage, 3, 4, 1)
```

* Network bytes per second tx

```
aliasByNode(scaleToSeconds(nonNegativeDerivative(collectd.$host.docker_stats.$app.$task.gauge.net.tx_bytes), 1), 3, 4, 1, 7)
```

* Network bytes per second rx

```
aliasByNode(scaleToSeconds(nonNegativeDerivative(collectd.$host.docker_stats.$app.$task.gauge.net.rx_bytes), 1), 3, 4, 1, 7)
```

* Network packets per second tx

```
aliasByNode(scaleToSeconds(nonNegativeDerivative(collectd.$host.docker_stats.$app.$task.gauge.net.tx_packets), 1), 3, 4, 1, 7)
```

* Network packets per second rx

```
aliasByNode(scaleToSeconds(nonNegativeDerivative(collectd.$host.docker_stats.$app.$task.gauge.net.rx_packets), 1), 3, 4, 1, 7)
```

## Running

Minimal command:

```
docker run -d -v /var/run/docker.sock:/var/run/docker.sock \
    -e GRAPHITE_HOST=<graphite host> -e COLLECTD_HOST=<colllectd host> \
    bobrik/collectd-docker
```

### Environment variables

* `COLLECTD_HOST` - host to use in metric name, defaults to `MESOS_HOST` if defined.
* `COLLECTD_INTERVAL` - metric update interval in seconds, defaults to `10`.
* `GRAPHITE_HOST` - host where carbon is listening for data.
* `GRAPHITE_PORT` - port where carbon is listening for data, `2003` by default.
* `GRAPHITE_PREFIX` - prefix for metrics in graphite, `collectd.` by default.
* `APP_LABEL_KEY` - container label to use for app name, `collectd_docker_app` by default.
* `APP_ENV_KEY` - container environment variable to use for app name, `COLLECTD_DOCKER_APP` by default.
* `TASK_LABEL_KEY` - container label to use for task name, `collectd_docker_task` by default.
* `TASK_ENV_KEY` - container environment variable to use for task name, `COLLECTD_DOCKER_TASK` by default.

Note that this docker image is very minimal and libc inside does not support
`search` directive in `/etc/resolv.conf`. You have to supply full hostname in
`GRAPHITE_HOST` that can be resolved with nameserver.

## License

MIT
