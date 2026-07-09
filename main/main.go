package main

import (
	"fmt"
	"time"

	"github.com/massarakhsh/lik/metric"
)

func main() {
	var counter metric.Count
	for n := 0; n < 1000; n++ {
		if n%2 == 0 || true {
			counter.Set(100)
		} else {
			counter.Set(0)
		}
		fmt.Printf("%d: %f\r\n", n, counter.Get())
		time.Sleep(time.Millisecond * 100)
	}
}
