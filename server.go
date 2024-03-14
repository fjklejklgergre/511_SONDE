package main

import (
  "encoding/gob"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net"
  "os"
)

type SystemInfo struct {
  Hostname   string `json:"hostname"`
  IP      string `json:"ip"`
  FreeRAM    float64 `json:"free_ram"`
  UsedRAM    float64 `json:"used_ram"`
  CPUPercentage float64 `json:"cpu_percentage"`
  FreeDiskSpace float64 `json:"free_disk_space"`
}

func collectSystemInfo() SystemInfo {
  hostInfo, _ := os.Hostname()
  ipAddress, _ := externalIP()

  vm, _ := getRAMInfo()
  cpuPercent, _ := getCPUInfo()
  diskStat, _ := getDiskInfo()

  return SystemInfo{
    Hostname:   hostInfo,
    IP:      ipAddress,
    FreeRAM:    float64(vm.Free) / 1024 / 1024 / 1024,
    UsedRAM:    float64(vm.Used) / 1024 / 1024 / 1024,
    CPUPercentage: cpuPercent,
    FreeDiskSpace: float64(diskStat.Free) / 1024 / 1024 / 1024,
  }
}

func getRAMInfo() (*MemoryInfo, error) {
  // ImplÃ©mentation pour obtenir les informations sur la RAM
  return nil, nil
}

func getCPUInfo() (float64, error) {
  // ImplÃ©mentation pour obtenir les informations sur le CPU
  return 0.0, nil
}

func getDiskInfo() (*DiskInfo, error) {
  // ImplÃ©mentation pour obtenir les informations sur le disque
  return nil, nil
}

type MemoryInfo struct {
  Total uint64
  Used uint64
  Free uint64
}

type DiskInfo struct {
  Total uint64
  Free uint64
}

func externalIP() (string, error) {
  // ImplÃ©mentation pour obtenir l'adresse IP externe
  return "", nil
}

func handleConnection(conn net.Conn) {
  defer conn.Close()

  decoder := gob.NewDecoder(conn)

  for {
    var info SystemInfo
    err := decoder.Decode(&info)
    if err != nil {
      fmt.Println("Error decoding data:", err)
      return
    }

    // Convertir SystemInfo en JSON
    jsonData, err := json.Marshal(info)
    if err != nil {
      fmt.Println("Error encoding data to JSON:", err)
      continue
    }

    // Actualiser les donnÃ©es JSON dans le fichier
    err = updateJSONInFile(jsonData, "/var/www/html/received_data.json")
    if err != nil {
      fmt.Println("Error updating JSON data in file:", err)
      continue
    }

    fmt.Println("Data updated in /var/www/html/received_data.json")
  }
}

func updateJSONInFile(newData []byte, filePath string) error {
  // Lire les donnÃ©es JSON existantes
  existingData, err := ioutil.ReadFile(filePath)
  if err != nil && !os.IsNotExist(err) {
    return err
  }

  // Ajouter un saut de ligne si le fichier n'est pas vide et n'existe pas
  if len(existingData) > 0 {
    newData = append([]byte("\n"), newData...)
  }

  // Ajouter les nouvelles donnÃ©es JSON Ã la fin des donnÃ©es existantes
  updatedData := append(existingData, newData...)

  // RÃ©Ã©crire le fichier avec les donnÃ©es mises Ã jour
  err = ioutil.WriteFile(filePath, updatedData, 0644)
  if err != nil {
    return err
  }

  return nil
}

func main() {
  listenAddr := ":8080"
  ln, err := net.Listen("tcp", listenAddr)
  if err != nil {
    fmt.Println("Error starting the server:", err)
    return
  }
  defer ln.Close()

  ipAddr := getOutboundIP()

  fmt.Printf("Server listening on %s:%s\n", ipAddr, listenAddr)

  for {
    conn, err := ln.Accept()
    if err != nil {
      fmt.Println("Error accepting connection:", err)
      continue
    }

    go handleConnection(conn)
  }
}

func getOutboundIP() string {
  conn, err := net.Dial("udp", "8.8.8.8:80")
  if err != nil {
    fmt.Println("Error getting outbound IP:", err)
    return ""
  }
  defer conn.Close()

  localAddr := conn.LocalAddr().(*net.UDPAddr)

  return localAddr.IP.String()
}
