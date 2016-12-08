package collector

import (
	"fmt"
	"io"
)

const collectdFloatGaugeTemplate = "PUTVAL %s/docker_stats-%s.%s/gauge-%s %d:%f\n"

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

func cpuPercentage(s Stats) float64 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(s.Stats.CPUStats.CPUUsage.TotalUsage) - float64(s.Stats.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta = float64(s.Stats.CPUStats.SystemCPUUsage) - float64(s.Stats.PreCPUStats.SystemCPUUsage)
	)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(s.Stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func (w CollectdWriter) writeInts(s Stats) error {
	metrics := map[string]float64{
		"cpu.user":       float64(s.Stats.CPUStats.CPUUsage.UsageInUsermode),
		"cpu.system":     float64(s.Stats.CPUStats.CPUUsage.UsageInKernelmode),
		"cpu.total":      float64(s.Stats.CPUStats.CPUUsage.TotalUsage),
		"cpu.percentage": cpuPercentage(s),

		"memory.limit": float64(s.Stats.MemoryStats.Limit),
		"memory.max":   float64(s.Stats.MemoryStats.MaxUsage),
		"memory.usage": float64(s.Stats.MemoryStats.Usage),

		"memory.active_anon":   float64(s.Stats.MemoryStats.Stats.TotalActiveAnon),
		"memory.active_file":   float64(s.Stats.MemoryStats.Stats.TotalActiveFile),
		"memory.cache":         float64(s.Stats.MemoryStats.Stats.TotalCache),
		"memory.inactive_anon": float64(s.Stats.MemoryStats.Stats.TotalInactiveAnon),
		"memory.inactive_file": float64(s.Stats.MemoryStats.Stats.TotalInactiveFile),
		"memory.mapped_file":   float64(s.Stats.MemoryStats.Stats.TotalMappedFile),
		"memory.pg_fault":      float64(s.Stats.MemoryStats.Stats.TotalPgfault),
		"memory.pg_in":         float64(s.Stats.MemoryStats.Stats.TotalPgpgin),
		"memory.pg_out":        float64(s.Stats.MemoryStats.Stats.TotalPgpgout),
		"memory.rss":           float64(s.Stats.MemoryStats.Stats.TotalRss),
		"memory.rss_huge":      float64(s.Stats.MemoryStats.Stats.TotalRssHuge),
		"memory.unevictable":   float64(s.Stats.MemoryStats.Stats.TotalUnevictable),
		"memory.writeback":     float64(s.Stats.MemoryStats.Stats.TotalWriteback),
	}

	for _, network := range s.Stats.Networks {
		metrics["net.rx_bytes"] += float64(network.RxBytes)
		metrics["net.rx_dropped"] += float64(network.RxDropped)
		metrics["net.rx_errors"] += float64(network.RxErrors)
		metrics["net.rx_packets"] += float64(network.RxPackets)

		metrics["net.tx_bytes"] += float64(network.TxBytes)
		metrics["net.tx_dropped"] += float64(network.TxDropped)
		metrics["net.tx_errors"] += float64(network.TxErrors)
		metrics["net.tx_packets"] += float64(network.TxPackets)
	}

	t := s.Stats.Read.Unix()

	for k, v := range metrics {
		err := w.writeFloat(s, k, t, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w CollectdWriter) writeFloat(s Stats, k string, t int64, v float64) error {
	msg := fmt.Sprintf(collectdFloatGaugeTemplate, w.host, s.App, s.Task, k, t, v)
	_, err := w.writer.Write([]byte(msg))
	return err
}
