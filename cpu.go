package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var (
	NumCPU         = runtime.NumCPU()
	CPUs, PrevCPUs [][10]int64
)

func init() {
	CPUs = make([][10]int64, NumCPU)
	PrevCPUs = make([][10]int64, NumCPU)
}

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
		displayUsage := int64(math.Floor(percentage * 10))
		if displayUsage == 10 {
			buf.WriteRune('X')
		} else {
			buf.WriteString(strconv.FormatInt(displayUsage, 10))
		}
		total += percentage
	}
	return fmt.Sprintf("%5.1f%% %s", total*100/float64(NumCPU), buf.String())
}

func moveCpuStatusToPrevious() {
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
	defer moveCpuStatusToPrevious()
	return fmt.Sprintf(" CPU%s ", calcCpuUsages())
}
