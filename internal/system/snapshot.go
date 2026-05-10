package system

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

type ProcessSortMode string

const (
	SortCPU    ProcessSortMode = "cpu"
	SortMemory ProcessSortMode = "memory"
)

type Snapshot struct {
	CPUPercent         float64       `json:"cpu_percent"`
	MemoryPercent      float64       `json:"memory_percent"`
	SwapPercent        float64       `json:"swap_percent"`
	DiskPercent        float64       `json:"disk_percent"`
	Load1              float64       `json:"load_1"`
	UptimeSeconds      uint64        `json:"uptime_seconds"`
	UptimeHuman        string        `json:"uptime_human"`
	TopCPUProcesses    []ProcessInfo `json:"top_cpu_processes"`
	TopMemoryProcesses []ProcessInfo `json:"top_memory_processes"`
}

type ProcessInfo struct {
	PID         int32   `json:"pid"`
	Name        string  `json:"name"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryMB    float32 `json:"memory_mb"`
	MemoryUsage float32 `json:"memory_usage_percent"`
}

func CollectSnapshot() (*Snapshot, error) {
	cpuValues, _ := cpu.Percent(500*time.Millisecond, false)
	memInfo, _ := mem.VirtualMemory()
	swapInfo, _ := mem.SwapMemory()
	diskInfo, _ := disk.Usage("/")
	loadInfo, _ := load.Avg()
	uptime, _ := host.Uptime()

	return &Snapshot{
		CPUPercent:         safePercent(cpuValues),
		MemoryPercent:      memInfo.UsedPercent,
		SwapPercent:        swapInfo.UsedPercent,
		DiskPercent:        diskInfo.UsedPercent,
		Load1:              loadInfo.Load1,
		UptimeSeconds:      uptime,
		UptimeHuman:        humanUptime(uptime),
		TopCPUProcesses:    collectTopProcesses(8, SortCPU),
		TopMemoryProcesses: collectTopProcesses(8, SortMemory),
	}, nil
}

func collectTopProcesses(limit int, sortMode ProcessSortMode) []ProcessInfo {
	processes, err := process.Processes()
	if err != nil {
		return nil
	}

	var result []ProcessInfo

	for _, p := range processes {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		cpuPercent, _ := p.CPUPercent()
		memPercent, _ := p.MemoryPercent()
		memInfo, _ := p.MemoryInfo()

		var memMB float32
		if memInfo != nil {
			memMB = float32(memInfo.RSS) / 1024 / 1024
		}

		result = append(result, ProcessInfo{
			PID:         p.Pid,
			Name:        name,
			CPUPercent:  cpuPercent,
			MemoryMB:    memMB,
			MemoryUsage: memPercent,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		switch sortMode {
		case SortMemory:
			return result[i].MemoryMB > result[j].MemoryMB
		default:
			return result[i].CPUPercent > result[j].CPUPercent
		}
	})

	if len(result) > limit {
		return result[:limit]
	}

	return result
}

func safePercent(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	return values[0]
}

func humanUptime(seconds uint64) string {
	d := seconds / 86400
	h := (seconds % 86400) / 3600
	m := (seconds % 3600) / 60

	if d > 0 {
		return fmt.Sprintf("%dd %dh %dm", d, h, m)
	}

	return fmt.Sprintf("%dh %dm", h, m)
}
