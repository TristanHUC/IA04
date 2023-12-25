package simulation

import (
	"fmt"
	"os"
)

type Counter struct {
	BeerCounterChan chan int
}

func (c *Counter) Run() {
	// Create a new file
	var nbAction int
	file, err := os.Create("beer_sales_log.txt")
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer file.Close()
	for {
		nbAction = <-c.BeerCounterChan

		// Write the current count and date to the file
		//currentTime := time.Now().Format(time.RFC3339)
		_, err := fmt.Fprintf(file, "%d\n", nbAction)
		if err != nil {
			fmt.Println("Error writing to file: ", err)
			return
		}

		// Flush any buffered data to the file
		err = file.Sync()
		if err != nil {
			fmt.Println("Error syncing file: ", err)
			return
		}
	}
}

func NewCounter() *Counter {
	return &Counter{
		BeerCounterChan: make(chan int, 10),
	}
}

func (c *Counter) GetChannelCounter() chan int {
	return c.BeerCounterChan
}
