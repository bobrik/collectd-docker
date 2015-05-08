package collector

import (
	"errors"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

const appLabel = "collectd_docker_app"
const taskLabel = "collectd_docker_task"

const appEnvPrefix = "COLLECTD_DOCKER_APP="
const taskEvnPrefix = "COLLECTD_DOCKER_TASK="

const defaultTask = "default"

// ErrNoNeedToMonitor is used to skip containers
// that shouldn't be monitored by collectd
var ErrNoNeedToMonitor = errors.New("container is not supposed to be monitored")

// Monitor is responsible for monitoring of a single container (task)
type Monitor struct {
	client   *docker.Client
	id       string
	app      string
	task     string
	interval int
}

// NewMonitor creates new monitor with specified docker client,
// container id and stat updating interval
func NewMonitor(c *docker.Client, id string, interval int) (*Monitor, error) {
	container, err := c.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	app := extractApp(container)
	if app == "" {
		return nil, ErrNoNeedToMonitor
	}

	task := extractTask(container)

	return &Monitor{
		client:   c,
		id:       container.ID,
		app:      app,
		task:     task,
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
				Task:  m.task,
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

func extractApp(c *docker.Container) string {
	return extractMetadata(c, appLabel, appEnvPrefix, "")
}

func extractTask(c *docker.Container) string {
	return extractMetadata(c, taskLabel, taskEvnPrefix, defaultTask)
}

func extractMetadata(c *docker.Container, label, envPrefix, missing string) string {
	if app, ok := c.Config.Labels[label]; ok {
		return app
	}

	for _, e := range c.Config.Env {
		if strings.HasPrefix(e, envPrefix) {
			return strings.TrimPrefix(e, envPrefix)
		}
	}

	return missing
}
