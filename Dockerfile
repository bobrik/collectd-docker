FROM debian:jessie

COPY . /go/src/github.com/bobrik/collectd-docker

RUN /go/src/github.com/bobrik/collectd-docker/docker/build.sh

ENTRYPOINT ["/run.sh"]
