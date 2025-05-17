package shutdown

import (
	"sync"
	"time"
)

type GracefulShutdown struct {
	ShutdownFn []func()

	config    config
	waitGroup *sync.WaitGroup
}

type config struct {
	GracePeriod   time.Duration
	CleanupPeriod time.Duration
}

func SetGracePeriodSeconds(s int64) func(c *config) {
	return func(c *config) {
		c.GracePeriod = time.Duration(s) * time.Second
	}
}

func SetCleanupPeriodSeconds(s int64) func(c *config) {
	return func(c *config) {
		c.CleanupPeriod = time.Duration(s) * time.Second
	}
}

func NewGracefulShutdown(fns []func(), cfgs ...func(c *config)) *GracefulShutdown {
	var (
		config = &config{}
	)

	for _, cfg := range cfgs {
		cfg(config)
	}

	return &GracefulShutdown{
		ShutdownFn: fns,
		config:     *config,
		waitGroup:  new(sync.WaitGroup),
	}
}

func (g *GracefulShutdown) Shutdown() {
	for _, fn := range g.ShutdownFn {
		g.waitGroup.Add(1)

		go func(fn func()) {
			defer g.waitGroup.Done()
			fn()
		}(fn)
	}

	g.waitGroup.Wait()
}
