package collector

import (
	"errors"
	"github.com/fsouza/go-dockerclient"
	"testing"
)

type fakeMonitorDockerClient struct {
	labels map[string]string
	env    []string
}

func (f fakeMonitorDockerClient) InspectContainer(id string) (*docker.Container, error) {
	return &docker.Container{
		Config: &docker.Config{
			Labels: f.labels,
			Env:    f.env,
		},
	}, nil
}

func (f fakeMonitorDockerClient) Stats(opts docker.StatsOptions) error {
	return errors.New("Stats() is not implemented for fake docker client")
}

func TestLabelExtraction(t *testing.T) {
	tests := map[*fakeMonitorDockerClient]struct {
		app  string
		task string
		err  error
	}{
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env:    []string{},
		}: {
			err: ErrNoNeedToMonitor,
		},

		// labels
		&fakeMonitorDockerClient{
			labels: map[string]string{
				appLabel: "myapp",
			},
			env: []string{},
		}: {
			app:  "myapp",
			task: defaultTask,
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{
				appLabel:  "myapp",
				taskLabel: "mytask",
			},
			env: []string{},
		}: {
			app:  "myapp",
			task: "mytask",
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{
				appLabel:  "my.app",
				taskLabel: "my.ta.sk",
			},
			env: []string{},
		}: {
			app:  "my_app",
			task: "my_ta_sk",
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{
				appLabel:          "my.app",
				taskLocationLabel: "some_label",
			},
			env: []string{},
		}: {
			app:  "my_app",
			task: defaultTask,
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{
				appLabel:          "my.app",
				taskLocationLabel: "some_label",
				"some_label":      "ho.ho.ho",
			},
			env: []string{},
		}: {
			app:  "my_app",
			task: "ho_ho_ho",
		},

		// env
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env: []string{
				appEnvPrefix + "myapp",
			},
		}: {
			app:  "myapp",
			task: defaultTask,
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env: []string{
				appEnvPrefix + "my.app",
				taskEnvPrefix + "my.ta.sk",
			},
		}: {
			app:  "my_app",
			task: "my_ta_sk",
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env: []string{
				appEnvPrefix + "myapp",
			},
		}: {
			app:  "myapp",
			task: defaultTask,
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env: []string{
				appEnvPrefix + "my.app",
				taskEnvLocationPrefix + "MESOS_TASK_ID",
				"MESOS_TASK_ID=topface_prod-test_app.c80a053f-f66f-11e4-a977-56847afe9799",
			},
		}: {
			app:  "my_app",
			task: "topface_prod-test_app_c80a053f-f66f-11e4-a977-56847afe9799",
		},
		&fakeMonitorDockerClient{
			labels: map[string]string{},
			env: []string{
				appEnvPrefix + "my.app",
				taskEnvLocationPrefix + "MESOS_TASK_ID",
				taskEnvLocationTrimPrefix + "topface_prod-test_app.",
				"MESOS_TASK_ID=topface_prod-test_app.c80a053f-f66f-11e4-a977-56847afe9799",
			},
		}: {
			app:  "my_app",
			task: "c80a053f-f66f-11e4-a977-56847afe9799",
		},
	}

	for c, e := range tests {
		m, err := NewMonitor(c, "", 1)
		if err != nil {
			if err != e.err {
				t.Errorf("expected error %q instead of %q for %#v", e.err, err, c)
			}

			continue
		} else {
			if e.err != nil {
				t.Errorf("expected error %q for %#v, got nothing", e.err, c)
				continue
			}
		}

		if m.app != e.app {
			t.Errorf("expected app %s got %s for %#v", e.app, m.app, c)
		}

		if m.task != e.task {
			t.Errorf("expected task %s got %s for %#v", e.task, m.task, c)
		}
	}

}
