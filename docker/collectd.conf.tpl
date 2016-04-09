Hostname "{{ .Env "COLLECTD_HOST" }}"

FQDNLookup false
Interval {{ .Env "COLLECTD_INTERVAL" }}
Timeout 2
ReadThreads 5

LoadPlugin write_graphite
<Plugin "write_graphite">
    <Carbon>
        Host "{{ .Env "GRAPHITE_HOST" }}"
        Port "{{ .Env "GRAPHITE_PORT" }}"
        Protocol "tcp"
        Prefix "{{ .Env "GRAPHITE_PREFIX" }}"
        StoreRates true
        EscapeCharacter "."
        AlwaysAppendDS false
        SeparateInstances true
    </Carbon>
</Plugin>

LoadPlugin exec
<Plugin exec>
  Exec "collectd-docker-collector" "/usr/bin/collectd-docker-collector" "-endpoint" "unix:///var/run/docker.sock" "-host" "{{ .Env "COLLECTD_HOST" }}" "-interval" "{{ .Env "COLLECTD_INTERVAL" }}"
</Plugin>
