package collector

import (
	"errors"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// ErrNoNeddToMonitor is used to skip containers
// that shouldn't be monitored by collectd
var ErrNoNeedToMonitor = errors.New("container is not supposed to be monitored")

// Monitor is responsible for monitoring of a single container
type Monitor struct {
	client   *docker.Client
	id       string
	app      string
	interval int
}

// NewMonitor creates new monitor with specified docker client,
// container id and stat updating interval
func NewMonitor(c *docker.Client, id string, interval int) (*Monitor, error) {
	container, err := c.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	app := extractCollectdApp(container)
	if app == "" {
		return nil, ErrNoNeedToMonitor
	}

	return &Monitor{
		client:   c,
		id:       container.ID,
		app:      app,
		interval: interval,
	}, nil
}

func (m *Monitor) handle(ch chan<- Stats) error {
	in := make(chan *docker.Stats)

	go func() {
		i := 0
		for s := range in {
			if i%m.interval != 0 {
				i++
				continue
			}

			ch <- Stats{
				App:   m.app,
				Stats: *s,
			}

			i++
		}
	}()

	return m.client.Stats(docker.StatsOptions{
		ID:    m.id,
		Stats: in,
	})
}

func extractCollectdApp(c *docker.Container) string {
	if app, ok := c.Config.Labels["collectd_app"]; ok {
		return app
	}

	for _, e := range c.Config.Env {
		if strings.HasPrefix(e, "COLLECTD_APP=") {
			return strings.TrimPrefix(e, "COLLECTD_APP=")
		}
	}

	return ""
}
