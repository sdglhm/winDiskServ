package main

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"net/http"
	"runtime"
)

type DiskUsage struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

func getDiskUsage(path string) (*DiskUsage, error) {
	diskStat, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}

	total := diskStat.Total
	free := diskStat.Free
	used := diskStat.Used
	usedPercent := diskStat.UsedPercent

	return &DiskUsage{
		Total:       total,
		Free:        free,
		Used:        used,
		UsedPercent: usedPercent,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>System Specifications</h1>")

	fmt.Fprintf(w, "<p>Go Version: %s</p>", runtime.Version())

	fmt.Fprintf(w, "<p>Operating System: %s</p>", runtime.GOOS)

	fmt.Fprintf(w, "<p>Number of CPU Cores: %d</p>", runtime.NumCPU())

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	fmt.Fprintf(w, "<p>Total Allocated Memory: %d bytes</p>", memStats.TotalAlloc)

	partitions, err := disk.Partitions(false)
	if err == nil {
		fmt.Fprintf(w, "<h2>Disk Space</h2>")
		for _, partition := range partitions {
			diskInfo, err := getDiskUsage(partition.Mountpoint)
			if err == nil {
				fmt.Fprintf(w, "<h3>%s</h3>", partition.Mountpoint)
				fmt.Fprintf(w, "<p>Total: %d bytes | %d GB</p>", diskInfo.Total, diskInfo.Total/1024/1024/1024)
				fmt.Fprintf(w, "<p>Free: %d bytes | %d GB</p>", diskInfo.Free, diskInfo.Free/1024/1024/1024)
				fmt.Fprintf(w, "<p>Used: %d bytes (%.2f%%)</p>", diskInfo.Used, diskInfo.UsedPercent)

				fmt.Fprintf(w, `<div style="border: 1px solid #ccc; border-radius: 5px; padding: 3px; width: 300px;">
					<div style="background-color: #4CAF50; height: 20px; width: %.2f%%; border-radius: 5px;"></div>
				</div>`, diskInfo.UsedPercent)
			}
		}
	}
}

func main() {
	port := 8090
	http.HandleFunc("/", handler)
	fmt.Printf("Server is running on http://0.0.0.0:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
