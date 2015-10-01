package main

import "cpuclock"

func main() {
	c := cpuclock.NewClock()
	go cpuclock.Routine(c)
	for c.Cnt < 100000 {

	}
	c.Stop = true
}
