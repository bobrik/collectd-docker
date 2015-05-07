package collector

import "github.com/fsouza/go-dockerclient"

// Stats represents singe stat from docker stats api for specific app
type Stats struct {
	App   string
	Stats docker.Stats
}
