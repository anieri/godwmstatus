package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	NumCPU         = runtime.NumCPU()
	CPUs, PrevCPUs [][10]int64
)

const (
	USER_I = iota
	NICE_I
	SYSTEM_I
	IDLE_I
	IOWAIT_I
	IRQ_I
	SOFT_IRQ_I
	STEAL_I
	GUEST_I
	GUEST_NICE_I
)

func procStat() string {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return " CPU - "
	}
	defer file.Close()

	var cpu string
	var user, nice, system, idle, iowait, irq, soft_irq, steal, guest, guest_nice int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(), "%s %d %d %d %d %d %d %d %d %d %d",
			&cpu, &user, &nice, &system, &idle, &iowait, &irq, &soft_irq, &steal, &guest, &guest_nice); err != nil {
			return " CPU - "
		}
		if strings.HasPrefix(cpu, "cpu") {
			cpuN := cpu[3:]
			if len(cpuN) == 0 {
				// First row
				continue
			}
			n, err := strconv.Atoi(cpuN)
			if err != nil {
				return " CPU - "
			}
			curr := &CPUs[n]
			curr[USER_I] = user
			curr[NICE_I] = nice
			curr[SYSTEM_I] = system
			curr[IDLE_I] = idle
			curr[IOWAIT_I] = iowait
			curr[IRQ_I] = irq
			curr[SOFT_IRQ_I] = soft_irq
			curr[STEAL_I] = steal
			curr[GUEST_I] = guest
			curr[GUEST_NICE_I] = guest_nice

			if n >= NumCPU-1 {
				break
			}
		}
	}
	defer moveCurrentStatToPrev()
	return fmt.Sprintf(" CPU%s ", calcCpuUsages())
}

// based on https://stackoverflow.com/questions/23367857/
func calcCpuUsages() string {
	var buf bytes.Buffer
	var total float64
	for i := 0; i < NumCPU; i++ {
		prev := PrevCPUs[i]
		curr := CPUs[i]

		prevIdle := prev[IDLE_I] + prev[IOWAIT_I]
		currIdle := curr[IDLE_I] + curr[IOWAIT_I]
		prevNonIdle := prev[USER_I] + prev[NICE_I] + prev[SYSTEM_I] + prev[IRQ_I] + prev[SOFT_IRQ_I] + prev[STEAL_I]
		currNonIdle := curr[USER_I] + curr[NICE_I] + curr[SYSTEM_I] + curr[IRQ_I] + curr[SOFT_IRQ_I] + curr[STEAL_I]
		prevTotal := prevIdle + prevNonIdle
		currTotal := currIdle + currNonIdle

		totald := currTotal - prevTotal
		idled := currIdle - prevIdle

		percentage := float64(totald-idled) / float64(totald)
		displayUsage := int(math.Floor(percentage * 10))
		if displayUsage == 10 {
			displayUsage = 9
		}
		total += percentage
		buf.WriteString(strconv.Itoa(displayUsage))
	}
	return fmt.Sprintf("%5.1f%% %s", total*100/float64(NumCPU), buf.String())
}

func moveCurrentStatToPrev() {
	for i := 0; i < NumCPU; i++ {
		prev := &PrevCPUs[i]
		curr := CPUs[i]
		prev[USER_I] = curr[USER_I]
		prev[NICE_I] = curr[NICE_I]
		prev[SYSTEM_I] = curr[SYSTEM_I]
		prev[IDLE_I] = curr[IDLE_I]
		prev[IOWAIT_I] = curr[IOWAIT_I]
		prev[IRQ_I] = curr[IRQ_I]
		prev[SOFT_IRQ_I] = curr[SOFT_IRQ_I]
		prev[STEAL_I] = curr[STEAL_I]
		prev[GUEST_I] = curr[GUEST_I]
		prev[GUEST_NICE_I] = curr[GUEST_NICE_I]
	}
}

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

const (
	RX = "↓"
	TX = "↑"
)

var (
	NetDevs = map[string]struct{}{
		"enp5s0:": {},
		"enp3s5:": {},
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

func main() {
	CPUs = make([][10]int64, NumCPU)
	PrevCPUs = make([][10]int64, NumCPU)
	for {
		var status = []string{
			procNetDev(),
			procStat(),
			procMeminfo(),
			time.Now().Local().Format(" Mon 02.01.2006 | 15:04:05"),
		}
		exec.Command("xsetroot", "-name", strings.Join(status, "|")).Run()
		var now = time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	}
}
