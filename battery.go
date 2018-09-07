package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

func batteryStatus(batteryDevice string) func() string {
	return func() string {
		out, err := exec.Command("upower", "-i", batteryDevice).CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return " BAT - "
		}

		var percentage, status, plugged string
		var parsedPercentage, parsedCharging bool
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			currentLine := scanner.Text()

			switch {
			case isPercentageLine(currentLine):
				percentage, status = readPercentage(currentLine)
				parsedPercentage = true
			case isChargingLine(currentLine):
				plugged = readCharging(currentLine)
				parsedCharging = true
			}

			if parsedPercentage && parsedCharging {
				if len(status) > 0 {
					return fmt.Sprintf(" %s %s %s ", status, plugged, percentage)
				}
				return fmt.Sprintf(" %s %s ", plugged, percentage)
			}
		}
		return " BAT - "
	}
}

var (
	percentageRegex = regexp.MustCompile(`percentage:`)
	chargingRegex   = regexp.MustCompile(`state:`)
)

func isPercentageLine(line string) bool {
	match := percentageRegex.FindString(line)
	return len(match) > 0
}

func isChargingLine(line string) bool {
	match := chargingRegex.FindString(line)
	return len(match) > 0
}

func readPercentage(line string) (string, string) {
	var percentage, void string
	if _, err := fmt.Sscanf(line, "%s %s", &void, &percentage); err != nil {
		fmt.Printf("percentage.parse.error: %v\n", err)
		return "-", ""
	}
	// strip percentage sign off the read value to parse int
	if i, err := strconv.ParseInt(percentage[:len(percentage)-1], 10, 32); err == nil {
		switch {
		case i < 10:
			return percentage, "CRITICAL"
		case i < 30:
			return percentage, "WARNING"
		}
	}
	return percentage, ""
}

func readCharging(line string) string {
	var chargingState, void string
	if _, err := fmt.Sscanf(line, "%s %s", &void, &chargingState); err != nil {
		fmt.Printf("plugged.parse.error: %v\n", err)
		return "BAT"
	}

	if chargingState == "charging" {
		return "PLG"
	}
	return "BAT"
}
