package collector

import (
	"log"
	"sync"

	"github.com/fsouza/go-dockerclient"
)

// Collector is responsible for discovering containers
// for monitoring and writing stats
type Collector struct {
	client     *docker.Client
	ch         chan Stats
	mutex      sync.Mutex
	registered map[string]struct{}
	interval   int
}

// NewCollector creates new Collector with specified docker client,
// collectd stats writer and stat updating interval
func NewCollector(client *docker.Client, w CollectdWriter, interval int) *Collector {
	ch := make(chan Stats)

	// TODO: this can be better, need to figure out how
	go func() {
		for s := range ch {
			w.Write(s)
		}
	}()

	return &Collector{
		client:     client,
		ch:         ch,
		mutex:      sync.Mutex{},
		registered: map[string]struct{}{},
		interval:   interval,
	}
}

// Run stats loop that discovers containers and runs
// monitoring tasks for them
func (c *Collector) Run(interval int) error {
	ch := make(chan *docker.APIEvents)
	err := c.client.AddEventListener(ch)
	if err != nil {
		return err
	}

	defer c.client.RemoveEventListener(ch)

	containers, err := c.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		go c.handle(container.ID)
	}

	for e := range ch {
		switch e.Status {
		case "start", "restart":
			go c.handle(e.ID)
		}
	}

	return nil
}

func (c *Collector) handle(id string) {
	m, err := NewMonitor(c.client, id, c.interval)
	if err != nil {
		if err == ErrNoNeedToMonitor {
			return
		}

		log.Printf("error handling %s: %s\n", id, err)

		return
	}

	go func() {
		if !c.register(id) {
			return
		}

		err := m.handle(c.ch)
		if err != nil {
			log.Printf("error handling container for app %s: %s\n", m.app, err)
		}

		c.unregister(id)
	}()
}

func (c *Collector) register(id string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.registered[id]; ok {
		return false
	}

	c.registered[id] = struct{}{}
	return true
}

func (c *Collector) unregister(id string) {
	c.mutex.Lock()
	delete(c.registered, id)
	c.mutex.Unlock()
}
