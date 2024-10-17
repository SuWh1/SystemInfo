package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Hostname     string `json:"hostname"`
	NumCPU       int    `json:"num_cpu"`
	GoVersion    string `json:"go_version"`

	Memory MemoryInfo `json:"memory"`
	CPU    CPUInfo    `json:"cpu"`
	Disk   DiskInfo   `json:"disk"`

	GPU []string `json:"gpu"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

type CPUInfo struct {
	ModelName string  `json:"model_name"`
	Cores     int32   `json:"cores"`
	Usage     float64 `json:"usage"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

func getGPUInfo() []string {
	switch runtime.GOOS {
	case "linux":
		return getGPUInfoLinux()
	case "windows":
		return getGPUInfoWindows()
	case "darwin":
		return getGPUInfoMacOS()
	default:
		return []string{"GPU information not available for this OS."}
	}
}

func getGPUInfoLinux() []string {
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		return []string{fmt.Sprintf("Error executing lspci: %v", err)}
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	var gpuInfo []string

	for _, line := range lines {
		if strings.Contains(line, "VGA") || strings.Contains(line, "3D") {
			gpuInfo = append(gpuInfo, line)
		}
	}

	if len(gpuInfo) == 0 {
		return []string{"No GPU information found."}
	}

	return gpuInfo
}

func getGPUInfoWindows() []string {
	cmd := exec.Command("powershell", "Get-WmiObject", "Win32_VideoController | Select-Object -Property Name,AdapterRAM | Format-Table -HideTableHeaders")
	output, err := cmd.Output()
	if err != nil {
		return []string{fmt.Sprintf("Error executing PowerShell command: %v", err)}
	}

	result := strings.TrimSpace(string(output))
	lines := strings.Split(result, "\n")

	var gpuInfo []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" { // Ensure that only non-empty lines are added
			gpuInfo = append(gpuInfo, line)
		}
	}

	if len(gpuInfo) == 0 {
		return []string{"No GPU information found."}
	}

	return gpuInfo
}

func getGPUInfoMacOS() []string {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		return []string{fmt.Sprintf("Error executing system_profiler: %v", err)}
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	var gpuInfo []string

	for _, line := range lines {
		if strings.Contains(line, "Chipset Model") || strings.Contains(line, "VRAM") {
			gpuInfo = append(gpuInfo, strings.TrimSpace(line))
		}
	}

	if len(gpuInfo) == 0 {
		return []string{"No GPU information found."}
	}

	return gpuInfo
}

func systemInfoHandler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	virtualMemory, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Info()
	cpuPercent, _ := cpu.Percent(0, false)
	diskStat, _ := disk.Usage("/")

	sysInfo := SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Hostname:     hostname,
		NumCPU:       runtime.NumCPU(),
		GoVersion:    runtime.Version(),
		Memory: MemoryInfo{
			Total:       virtualMemory.Total / 1024 / 1024,
			Free:        virtualMemory.Free / 1024 / 1024,
			Used:        virtualMemory.Used / 1024 / 1024,
			UsedPercent: virtualMemory.UsedPercent,
		},
		CPU: CPUInfo{
			ModelName: cpuInfo[0].ModelName,
			Cores:     cpuInfo[0].Cores,
			Usage:     cpuPercent[0],
		},
		Disk: DiskInfo{
			Total:       diskStat.Total / 1024 / 1024 / 1024,
			Free:        diskStat.Free / 1024 / 1024 / 1024,
			Used:        diskStat.Used / 1024 / 1024 / 1024,
			UsedPercent: diskStat.UsedPercent,
		},
		GPU: getGPUInfo(),
	}

	tmpl, err := template.ParseFiles("templates/system.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, sysInfo); err != nil {
		fmt.Printf("Error rendering template: %v\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/main.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/systeminfo", systemInfoHandler)
	http.ListenAndServe(":8080", nil)
}
