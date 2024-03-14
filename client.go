package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/host"
)

type SystemInfo struct {
	Hostname      string
	IP            string
	FreeRAM       float64
	UsedRAM       float64
	CPUPercentage float64
	FreeDiskSpace float64
}

func collectSystemInfo() SystemInfo {
	// Collect hostname and IP information
	hostInfo, _ := host.Info()
	interfaces, _ := net.Interfaces()

	var ipAddress string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && !strings.Contains(iface.Name, "lo") {
			addrs, _ := iface.Addrs()
			if len(addrs) > 0 {
				ipAddress = strings.Split(addrs[0].String(), "/")[0]
				break
			}
		}
	}

	// Collect RAM information
	vm, _ := mem.VirtualMemory()

	// Collect CPU information
	cpuPercent, _ := cpu.Percent(0, false)

	// Collect Disk information
	diskStat, _ := disk.Usage("/")

	return SystemInfo{
		Hostname:      hostInfo.Hostname,
		IP:            ipAddress,
		FreeRAM:       float64(vm.Free) / 1024 / 1024 / 1024, // convert to GB
		UsedRAM:       float64(vm.Used) / 1024 / 1024 / 1024, // convert to GB
		CPUPercentage: cpuPercent[0],
		FreeDiskSpace: float64(diskStat.Free) / 1024 / 1024 / 1024, // convert to GB
	}
}

func main() {
	fmt.Print("Enter the server IP address: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	serverAddr := scanner.Text()

	conn, err := net.Dial("tcp", serverAddr+":8080")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)

	for {
		info := collectSystemInfo()
		err := encoder.Encode(info)
		if err != nil {
			fmt.Println("Error encoding and sending data:", err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}
