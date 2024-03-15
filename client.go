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

// Structure pour stocker les informations système
type SystemInfo struct {
	Hostname      string  // Nom d'hôte
	IP            string  // Adresse IP
	FreeRAM       float64 // RAM libre (en Go)
	UsedRAM       float64 // RAM utilisée (en Go)
	CPUPercentage float64 // Pourcentage d'utilisation du CPU
	FreeDiskSpace float64 // Espace disque libre (en Go)
}

// Fonction pour collecter les informations système
func collectSystemInfo() SystemInfo {
	// Collecte des informations sur le nom d'hôte et l'adresse IP
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

	// Collecte des informations sur la RAM
	vm, _ := mem.VirtualMemory()

	// Collecte des informations sur le CPU
	cpuPercent, _ := cpu.Percent(0, false)

	// Collecte des informations sur le disque
	diskStat, _ := disk.Usage("/")

	return SystemInfo{
		Hostname:      hostInfo.Hostname,
		IP:            ipAddress,
		FreeRAM:       float64(vm.Free) / 1024 / 1024 / 1024, // Conversion en Go
		UsedRAM:       float64(vm.Used) / 1024 / 1024 / 1024, // Conversion en Go
		CPUPercentage: cpuPercent[0],
		FreeDiskSpace: float64(diskStat.Free) / 1024 / 1024 / 1024, // Conversion en Go
	}
}

func main() {
	fmt.Print("Entrez l'adresse IP du serveur : ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	serverAddr := scanner.Text()

	conn, err := net.Dial("tcp", serverAddr+":8080")
	if err != nil {
		fmt.Println("Erreur lors de la connexion au serveur :", err)
		return
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)

	for {
		info := collectSystemInfo()
		err := encoder.Encode(info)
		if err != nil {
			fmt.Println("Erreur lors de l'encodage et de l'envoi des données :", err)
			return
		}

		time.Sleep(5 * time.Second) // Attente de 5 secondes avant de collecter à nouveau les informations
	}
}
