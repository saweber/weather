package main

import (
	"fmt"
	"time"
)

func main() {
	duration, err := time.ParseDuration("24h")
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		// intentionally fail to make this problem obvious
		panic("unable to setup scheduled task due to bad duration")
	}

	// run immediately to get most recent NOAA storm reports
	getNOAAStormReports()

	// schedule collection of future NOAA storm reports
	scheduleTask(duration, getNOAAStormReports)
}

func scheduleTask(duration time.Duration, task func()) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			task()
		}
	}
}
