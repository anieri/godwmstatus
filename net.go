package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	RX = "↓"
	TX = "↑"
)

var (
	NetDevs = map[string]struct{}{
		"enp2s0:": {},
		"wlp3s0:": {},
	}
	prevRX, prevTX int64
)

func procNetDev() string {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return fmt.Sprintf("%s - %s - ", RX, TX)
	}
	defer file.Close()

	var dev string
	var rx, tx, currRX, currTX, void int64

	var scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		if _, err = fmt.Sscanf(scanner.Text(), "%s %d %d %d %d %d %d %d %d %d",
			&dev, &rx, &void, &void, &void, &void, &void, &void, &void, &tx); err != nil {
			continue
		}
		if _, ok := NetDevs[dev]; ok {
			currRX += rx
			currTX += tx
		}
	}

	defer func() { prevRX, prevTX = currRX, currTX }()
	return fmt.Sprintf("%s  %s  ", fixed(RX, currRX-prevRX), fixed(TX, currTX-prevTX))
}

func fixed(prefix string, rate int64) string {
	if rate < 0 {
		return fmt.Sprintf("%s -", prefix)
	}

	var decDigit int64
	suffix := "B"

	switch {
	case rate >= (1000 * 1024):
		decDigit = (rate / 1024 / 102) % 10
		rate /= (1024 * 1024)
		suffix = "M"
	case rate >= 1000:
		decDigit = (rate / 102) % 10
		rate /= 1024
		suffix = "K"
	}

	if rate >= 100 {
		return fmt.Sprintf("%s%4d%s", prefix, rate, suffix)
	}
	return fmt.Sprintf("%s%2d.%1d%s", prefix, rate, decDigit, suffix)
}
