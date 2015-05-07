package collector

import (
	"fmt"
	"io"
)

const collectdIntGaugeTemplate = "PUTVAL %s/docker_stats-%s/gauge-%s %d:%d\n"

// CollectdWriter is responsible for writing data
// to wrapped writer in collectd exec plugin format
type CollectdWriter struct {
	host     string
	writer   io.Writer
	interval int
}

// NewCollectdWriter creates new CollectdWriter
// with specified hostname and writer
func NewCollectdWriter(host string, writer io.Writer) CollectdWriter {
	return CollectdWriter{
		host:   host,
		writer: writer,
	}
}

func (w CollectdWriter) Write(s Stats) error {

	return w.writeInts(s)
}

func (w CollectdWriter) writeInts(s Stats) error {
	metrics := map[string]int64{
		"cpu.user":   s.Stats.CPUStats.CPUUsage.UsageInUsermode,
		"cpu.system": s.Stats.CPUStats.CPUUsage.UsageInKernelmode,
		"cpu.total":  s.Stats.CPUStats.CPUUsage.TotalUsage,

		"memory.limit": s.Stats.MemoryStats.Limit,
		"memory.max":   s.Stats.MemoryStats.MaxUsage,
		"memory.usage": s.Stats.MemoryStats.Usage,

		"memory.active_anon":   s.Stats.MemoryStats.Stats.TotalActiveAnon,
		"memory.active_file":   s.Stats.MemoryStats.Stats.TotalActiveFile,
		"memory.cache":         s.Stats.MemoryStats.Stats.TotalCache,
		"memory.inactive_anon": s.Stats.MemoryStats.Stats.TotalInactiveAnon,
		"memory.inactive_file": s.Stats.MemoryStats.Stats.TotalInactiveFile,
		"memory.mapped_file":   s.Stats.MemoryStats.Stats.TotalMappedFile,
		"memory.pg_fault":      s.Stats.MemoryStats.Stats.TotalPgfault,
		"memory.pg_in":         s.Stats.MemoryStats.Stats.TotalPgpgin,
		"memory.pg_out":        s.Stats.MemoryStats.Stats.TotalPgpgout,
		"memory.rss":           s.Stats.MemoryStats.Stats.TotalRss,
		"memory.rss_huge":      s.Stats.MemoryStats.Stats.TotalRssHuge,
		"memory.unevictable":   s.Stats.MemoryStats.Stats.TotalUnevictable,
		"memory.writeback":     s.Stats.MemoryStats.Stats.TotalWriteback,

		"net.rx_bytes":   s.Stats.Network.RxBytes,
		"net.rx_dropped": s.Stats.Network.RxDropped,
		"net.rx_errors":  s.Stats.Network.RxErrors,
		"net.rx_packets": s.Stats.Network.RxPackets,
		"net.tx_bytes":   s.Stats.Network.TxBytes,
		"net.tx_dropped": s.Stats.Network.TxDropped,
		"net.tx_errors":  s.Stats.Network.TxErrors,
		"net.tx_packets": s.Stats.Network.TxPackets,
	}

	t := s.Stats.Read.Unix()

	for k, v := range metrics {
		err := w.writeInt(s.App, k, t, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w CollectdWriter) writeInt(a, k string, t, v int64) error {
	msg := fmt.Sprintf(collectdIntGaugeTemplate, w.host, a, k, t, v)
	_, err := w.writer.Write([]byte(msg))
	return err
}
