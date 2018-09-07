package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func batteryStatus() string {
	out, err := exec.Command("sh", "-c", "upower -i /org/freedesktop/UPower/devices/battery_BAT0 | grep percentage").CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return " BAT - "
	}

	var percentage, void string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(), "%s %s", &void, &percentage); err != nil {
			fmt.Printf("Error: %v\n", err)
			return " BAT - "
		}
		if i, err := strconv.ParseInt(percentage[:len(percentage)-1], 10, 32); err == nil {
			switch {
			case i < 10:
				return fmt.Sprintf(" CRITICAL BAT %s ", percentage)
			case i < 30:
				return fmt.Sprintf(" WARNING BAT %s ", percentage)
			}
		}
		return fmt.Sprintf(" BAT %s ", percentage)
	}
	return " BAT - "
}
