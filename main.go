package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func getGPUInfoLinux() string {
	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing lspci: %v\n", err)
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if strings.Contains(line, "VGA") {
			return fmt.Sprintf("GPU Info: %s\n", line)
		}
	}
	return "No GPU information found."
}

func getGPUInfoNvidia() string {
	cmd := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing nvidia-smi: %v\n", err)
	}

	return fmt.Sprintf("NVIDIA GPU Info:\n%s\n", string(output))
}

func getGPUInfoWindows() string {
	cmd := exec.Command("powershell", "Get-WmiObject", "Win32_VideoController | Select-Object -Property Name,AdapterRAM | Format-Table -HideTableHeaders")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing PowerShell command: %v", err)
	}

	result := strings.TrimSpace(string(output))
	return fmt.Sprintf("Windows GPU Info:\n%s\n", result)
}

func getGPUInfoMacOS() string {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error executing system_profiler: %v", err)
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	var gpuInfo []string
	for _, line := range lines {
		if strings.Contains(line, "Chipset Model") || strings.Contains(line, "VRAM") {
			gpuInfo = append(gpuInfo, strings.TrimSpace(line))
		}
	}
	return strings.Join(gpuInfo, "\n")
}

func main() {
	fmt.Println("==================")
	fmt.Println("System Information")
	fmt.Println("==================")
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)

	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error fetching hostname: %v\n", err)
	} else {
		fmt.Printf("Device name: %s\n", host)
	}

	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())

	fmt.Printf("Go Version: %s\n", runtime.Version())

	virtualMem, _ := mem.VirtualMemory()
	fmt.Println("\n==================")
	fmt.Printf("Memory Information\n")
	fmt.Println("==================")
	fmt.Printf("Total: %v MB\n", virtualMem.Total/1024/1024)
	fmt.Printf("Free: %v MB\n", virtualMem.Free/1024/1024)
	fmt.Printf("Used: %v MB\n", virtualMem.Used/1024/1024)
	fmt.Printf("Used Percent: %0.3f%%\n", virtualMem.UsedPercent)

	fmt.Println("\n==================")
	fmt.Printf("CPU Information\n")
	fmt.Println("==================")
	cpuInfo, _ := cpu.Info()
	for _, info := range cpuInfo {
		fmt.Printf("Model: %s\n", info.ModelName)
		fmt.Printf("Cores: %d\n", info.Cores)
	}
	cpuPercent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("CPU Usage: %.2f%%\n\n", cpuPercent[0])

	fmt.Println("==================")
	fmt.Printf("Disk Information\n")
	fmt.Println("==================")
	diskStat, _ := disk.Usage("/")
	fmt.Printf("Total: %v GB\n", diskStat.Total/1024/1024/1024)
	fmt.Printf("Free: %v GB\n", diskStat.Free/1024/1024/1024)
	fmt.Printf("Used: %v GB\n", diskStat.Used/1024/1024/1024)
	fmt.Printf("Used Percent: %0.3f%%\n", diskStat.UsedPercent)

	fmt.Println("\n==================")
	fmt.Printf("GPU Information\n")
	fmt.Println("==================")
	switch runtime.GOOS {
	case "linux":
		nvidiaInfo := getGPUInfoNvidia()
		if strings.Contains(nvidiaInfo, "Error") {
			fmt.Println(getGPUInfoLinux())
		} else {
			fmt.Println(nvidiaInfo)
		}
	case "windows":
		fmt.Println(getGPUInfoWindows())
	case "darwin": // macOS
		fmt.Println(getGPUInfoMacOS())
	default:
		fmt.Println("GPU information not available for this OS.")
	}
}
