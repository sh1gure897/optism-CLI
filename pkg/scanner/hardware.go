package scanner

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	OS          string
	Arch        string
	CPUName     string
	CPUCores    int
	TotalRAM_MB uint64
}

func ScanHardware() (*SystemInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %v", err)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get cpu info: %v", err)
	}

	cpuName := "Unknown CPU"
	if len(cpuInfo) > 0 {
		cpuName = cpuInfo[0].ModelName
	}

	return &SystemInfo{
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
		CPUName:     cpuName,
		CPUCores:    runtime.NumCPU(),
		TotalRAM_MB: v.Total / 1024 / 1024,
	}, nil
}
