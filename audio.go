package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var (
	volumeRegex = regexp.MustCompile(`\[[0-9]+%\]`)
)

func amixerVolume() string {
	out, err := exec.Command("sh", "-c", "amixer sget Master").CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return " VOL - "
	}

	var totalVolume, nVolume int
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		match := volumeRegex.FindString(scanner.Text())
		if len(match) <= 3 {
			continue
		}
		vol, err := strconv.Atoi(match[1 : len(match)-2])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		totalVolume += vol
		nVolume++
	}
	if nVolume == 0 {
		nVolume = 1
	}
	return fmt.Sprintf(" VOL %d%% ", totalVolume/nVolume)
}
