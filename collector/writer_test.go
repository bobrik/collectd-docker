package collector

import (
	"github.com/fsouza/go-dockerclient"
	"testing"
)

type CPUUsage struct {
	PercpuUsage       []uint64 `json:"percpu_usage,omitempty" yaml:"percpu_usage,omitempty"`
	UsageInUsermode   uint64   `json:"usage_in_usermode,omitempty" yaml:"usage_in_usermode,omitempty"`
	TotalUsage        uint64   `json:"total_usage,omitempty" yaml:"total_usage,omitempty"`
	UsageInKernelmode uint64   `json:"usage_in_kernelmode,omitempty" yaml:"usage_in_kernelmode,omitempty"`
}

func TestCpuPercentaget(t *testing.T) {
	tests := []Stats{
		Stats{
			App:   "app-test",
			Task:  "task-test",
			Stats: docker.Stats{},
		},
		Stats{
			App:  "app-test",
			Task: "task-test",
			Stats: docker.Stats{
				CPUStats: docker.CPUStats{
					SystemCPUUsage: uint64(31398238640000000),
					CPUUsage: CPUUsage{
						TotalUsage: 2594015077802,
						PercpuUsage: []uint64{
							323546750673,
							324564397608,
							317061012401,
							316836458442,
							281575177161,
							275684143003,
							43863240551,
							44670161231,
							43724548993,
							47349938842,
							43608875064,
							44060577573,
							23158407180,
							18979754536,
							95987831450,
							119531282932,
							101243729356,
							106668938153,
							3805746869,
							5166565480,
							2801432984,
							2098161188,
							3225196915,
							4802749217,
						},
					},
				},
				PreCPUStats: docker.CPUStats{
					SystemCPUUsage: uint64(0),
					CPUUsage: CPUUsage{
						TotalUsage: 0,
					},
				},
			},
		},
		Stats{
			App:  "app-test",
			Task: "task-test",
			Stats: docker.Stats{
				CPUStats: docker.CPUStats{
					SystemCPUUsage: 31398262460000000,
					CPUUsage: CPUUsage{
						TotalUsage: 2000,
						PercpuUsage: []uint64{
							323546750673,
							324564397608,
							317061012401,
							316836458442,
							281575177161,
							275684143003,
							43863240551,
							44670161231,
							43724548993,
							47349938842,
							43608875064,
							44060577573,
							23158407180,
							18979754536,
							95987831450,
							119531282932,
							101243729356,
							106668938153,
							3805746869,
							5166565480,
							2801432984,
							2098161188,
							3225196915,
							4802749217,
						},
					},
				},
				PreCPUStats: docker.CPUStats{
					SystemCPUUsage: 31398238640000000,
					CPUUsage: CPUUsage{
						TotalUsage: 2594015077802,
					},
				},
			},
		},
	}

	for c, e := range tests {
		m := cpuPercentage(e)
		if m < 0.0 || m > 100.0 {
			t.Errorf("[%d] expected cpu percentage within bounds but got %f", c, m)
		}
	}
}
