package main

import (
	"bufio"
	"fmt"
	"os"
)

func procMeminfo() string {
	var file, err = os.Open("/proc/meminfo")
	if err != nil {
		return " MEM - "
	}
	defer file.Close()

	var total, used, done int

	scanner := bufio.NewScanner(file)
	for done != 15 && scanner.Scan() {
		var prop, val = "", 0
		if _, err = fmt.Sscanf(scanner.Text(), "%s %d", &prop, &val); err != nil {
			return " MEM - "
		}
		switch prop {
		case "MemTotal:":
			total = val
			used += val
			done |= 1
		case "MemFree:":
			used -= val
			done |= 2
		case "Buffers:":
			used -= val
			done |= 4
		case "Cached:":
			used -= val
			done |= 8
		}
	}
	percentage := used * 100 / total
	return fmt.Sprintf(" MEM%3d%% ", percentage)
}
