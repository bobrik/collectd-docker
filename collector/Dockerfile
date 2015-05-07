FROM alpine:3.1

RUN apk --update add collectd py-pip && \
    pip install envtpl

COPY ./docker/collectd.conf.tpl /etc/collectd/collectd.conf.tpl

COPY ./docker/run.sh /run.sh

COPY . /go/src/github.com/bobrik/collectd-docker/collector

RUN apk --update add go git && \
    GOPATH=/go go get github.com/bobrik/collectd-docker/collector/... && \
    apk del go git && \
    mv /go/bin/collector /collector && \
    rm -rf /go && \
    chmod 6755 /collector

ENTRYPOINT ["/run.sh"]
