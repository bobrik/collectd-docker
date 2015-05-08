package collector

import "github.com/fsouza/go-dockerclient"

// Stats represents singe stat from docker stats api for specific task
type Stats struct {
	App   string
	Task  string
	Stats docker.Stats
}
