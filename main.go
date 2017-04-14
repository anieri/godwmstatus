package main

import (
	"bufio"
	"log"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	NumCPU         = runtime.NumCPU()
	CPUs, PrevCpus [][10]int64
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
		return "CPU -"
	}
	defer file.Close()

	var cpu string
	var user, nice, system, idle, iowait, irq, soft_irq, steal, guest, guest_nice int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if _, err := fmt.Sscanf(scanner.Text(), "%s %d %d %d %d %d %d %d %d %d %d",
			&cpu, &user, &nice, &system, &idle, &iowait, &irq, &soft_irq, &steal, &guest, &guest_nice); err != nil {
			return "CPU"
		}
		if strings.HasPrefix(cpu, "cpu") {
			cpuN := cpu[3:]
			if len(cpuN) == 0 {
				// First row
				continue
			}
			n, err := strconv.Atoi(cpuN)
			if err != nil {
				return "CPU"
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
	for i, o := range CPUs {
		log.Printf("CPU %d: %+v\n", i, o)
	}

	return fmt.Sprintf("CPU")
}

func main() {
	CPUs = make([][10]int64, NumCPU)
	PrevCpus = make([][10]int64, NumCPU)
	for {
		var status = []string{
			procStat(),
			time.Now().Local().Format("Mon 02 Jan 2006 | 15:04:05"),
		}
		exec.Command("xsetroot", "-name", strings.Join(status, " ")).Run()
		var now = time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	}
}
