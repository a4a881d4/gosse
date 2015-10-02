package main

import "cpuclock"
import "time"
import "fmt"

func main() {
	c := cpuclock.NewClock()
	go cpuclock.Routine(c)
	for c.Cnt < 100000 {
		time.Sleep(time.Second * 1)
		r, s := c.Run()
		fmt.Println("time diff", r+s-time.Now().UnixNano(), c.Cnt)
	}
	c.Stop = true
}
