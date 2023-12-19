package simulation

import "fmt"

type Counter struct {
	BeerCounterChan chan bool
	Counter         int
}

func (c *Counter) Run() {
	fmt.Println("Counter started")
	for {
		fmt.Println("Counter waiting")
		_ = <-c.BeerCounterChan
		c.Counter++
		if c.Counter%10 == 0 {
			println(c.Counter)
		}
	}
}

func (c *Counter) GetCounter() int {
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
