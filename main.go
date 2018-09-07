package main

import (
	"os/exec"
	"strings"
	"time"
)

func throttle(rate int, fn func() string) func() string {
	previousResult := ""
	currentI := 0
	return func() string {
		if (currentI % rate) == 0 {
			previousResult = fn()
		}
		currentI++
		return previousResult
	}
}

func main() {
	_amixerVolume := throttle(10, amixerVolume)
	_batteryStatus := throttle(90, batteryStatus)

	for {
		var status = []string{
			procNetDev(),
			procStat(),
			procMeminfo(),
			_amixerVolume(),
			_batteryStatus(),
			time.Now().Local().Format(" Mon 02.01.2006 | 15:04:05"),
		}
		exec.Command("xsetroot", "-name", strings.Join(status, "|")).Run()
		var now = time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	}
}
