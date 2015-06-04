Hostname "{{ COLLECTD_HOST }}"

FQDNLookup false
Interval {{ COLLECTD_INTERVAL | default(10) }}
Timeout 2
ReadThreads 5

LoadPlugin write_graphite
<Plugin "write_graphite">
    <Carbon>
        Host "{{ GRAPHITE_HOST }}"
        Port "{{ GRAPHITE_PORT | default("2003") }}"
        Protocol "tcp"
        Prefix "{{ GRAPHITE_PREFIX | default("collectd.") }}"
        StoreRates true
        EscapeCharacter "."
        AlwaysAppendDS false
        SeparateInstances true
    </Carbon>
</Plugin>

LoadPlugin exec
<Plugin exec>
  Exec "nobody" "/collector" "-endpoint" "unix:///var/run/docker.sock" "-host" {{ COLLECTD_HOST }} "-interval" "{{ COLLECTD_INTERVAL | default(10) }}"
</Plugin>
