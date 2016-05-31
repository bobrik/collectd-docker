package collector

import (
	"fmt"
	"io"
)

const collectdIntGaugeTemplate = "PUTVAL %s/docker_stats-%s.%s/gauge-%s %d:%d\n"
const collectdFloatPercentageTemplate = "PUTVAL %s/docker_stats-%s.%s/percent-%s %d:%f\n"

// CollectdWriter is responsible for writing data
// to wrapped writer in collectd exec plugin format
type CollectdWriter struct {
	host      string
	writer    io.Writer
	interval  int
}

// NewCollectdWriter creates new CollectdWriter
// with specified hostname and writer
func NewCollectdWriter(host string, writer io.Writer) CollectdWriter {
	return CollectdWriter{
		host:   host,
		writer: writer,
	}
}

func (w *CollectdWriter) Write(s Stats) error {
	if err := w.writeInts(s); err != nil {
		return err
	}

	if err := w.writeDerived(s); err != nil {
		return err
	}

	return nil
}

func (w CollectdWriter) writeInts(s Stats) error {
	metrics := map[string]uint64{
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
	}

	for _, network := range s.Stats.Networks {
		metrics["net.rx_bytes"] += network.RxBytes
		metrics["net.rx_dropped"] += network.RxDropped
		metrics["net.rx_errors"] += network.RxErrors
		metrics["net.rx_packets"] += network.RxPackets

		metrics["net.tx_bytes"] += network.TxBytes
		metrics["net.tx_dropped"] += network.TxDropped
		metrics["net.tx_errors"] += network.TxErrors
		metrics["net.tx_packets"] += network.TxPackets
	}

	t := s.Stats.Read.Unix()

	for k, v := range metrics {
		err := w.writeInt(s, k, t, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w CollectdWriter) writeDerived(s Stats) error {
	cpu := w.calculateCPUPercent(s)

	t := s.Stats.Read.Unix()

	return w.writeFloat(s, "cpu.percentage", t, cpu)
}

func (w CollectdWriter) writeInt(s Stats, k string, t int64, v uint64) error {
	msg := fmt.Sprintf(collectdIntGaugeTemplate, w.host, s.App, s.Task, k, t, v)
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w CollectdWriter) writeFloat(s Stats, k string, t int64, v float64) error {
	msg := fmt.Sprintf(collectdFloatPercentageTemplate, w.host, s.App, s.Task, k, t, v)
	_, err := w.writer.Write([]byte(msg))
	return err
}

// calculate CPU Percentage, taken from docker stats client:
// https://github.com/docker/docker/blob/master/api/client/stats_helpers.go#L199-L212
func (w CollectdWriter) calculateCPUPercent(s Stats) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(s.Stats.CPUStats.CPUUsage.TotalUsage) - float64(s.PrevStats.CPUStats.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(s.Stats.CPUStats.SystemCPUUsage) - float64(s.PrevStats.CPUStats.SystemCPUUsage)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(s.Stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return cpuPercent
}
