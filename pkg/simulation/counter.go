package simulation

type Counter struct {
	BeerCounterChan chan bool
	Counter         int
	running         bool
}

func (c *Counter) Run() {
	c.running = true
	for c.running {
		_ = <-c.BeerCounterChan
		c.Counter++
	}
}

func (c *Counter) GetCount() int {
	return c.Counter
}

func NewCounter() *Counter {
	return &Counter{
		BeerCounterChan: make(chan bool, 10),
		Counter:         0,
	}
}

func (c *Counter) GetChannelCounter() chan bool {
	return c.BeerCounterChan
}
