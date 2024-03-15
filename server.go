package main

import (
  "encoding/gob"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net"
  "os"
)

// Structure pour stocker les informations système
type SystemInfo struct {
  Hostname      string  `json:"hostname"`       // Nom d'hôte de la machine
  IP            string  `json:"ip"`             // Adresse IP de la machine
  FreeRAM       float64 `json:"free_ram"`       // RAM libre en Go
  UsedRAM       float64 `json:"used_ram"`       // RAM utilisée en Go
  CPUPercentage float64 `json:"cpu_percentage"` // Pourcentage d'utilisation du CPU
  FreeDiskSpace float64 `json:"free_disk_space"`// Espace disque libre en Go
}

// Fonction pour collecter les informations système
func collectSystemInfo() SystemInfo {
  hostInfo, _ := os.Hostname()
  ipAddress, _ := externalIP()

  vm, _ := getRAMInfo()      // Obtenir les informations sur la RAM
  cpuPercent, _ := getCPUInfo() // Obtenir les informations sur le CPU
  diskStat, _ := getDiskInfo()  // Obtenir les informations sur le disque

  return SystemInfo{
    Hostname:      hostInfo,
    IP:            ipAddress,
    FreeRAM:       float64(vm.Free) / 1024 / 1024 / 1024,
    UsedRAM:       float64(vm.Used) / 1024 / 1024 / 1024,
    CPUPercentage: cpuPercent,
    FreeDiskSpace: float64(diskStat.Free) / 1024 / 1024 / 1024,
  }
}

// Fonction pour obtenir les informations sur la RAM
func getRAMInfo() (*MemoryInfo, error) {
  // Implémentation pour obtenir les informations sur la RAM
  return nil, nil
}

// Fonction pour obtenir les informations sur le CPU
func getCPUInfo() (float64, error) {
  // Implémentation pour obtenir les informations sur le CPU
  return 0.0, nil
}

// Fonction pour obtenir les informations sur le disque
func getDiskInfo() (*DiskInfo, error) {
  // Implémentation pour obtenir les informations sur le disque
  return nil, nil
}

// Structure pour stocker les informations sur la mémoire
type MemoryInfo struct {
  Total uint64
  Used  uint64
  Free  uint64
}

// Structure pour stocker les informations sur le disque
type DiskInfo struct {
  Total uint64
  Free  uint64
}

// Fonction pour obtenir l'adresse IP externe
func externalIP() (string, error) {
  // Implémentation pour obtenir l'adresse IP externe
  return "", nil
}

// Fonction pour gérer les connexions entrantes
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

    // Actualiser les données JSON dans le fichier
    err = updateJSONInFile(jsonData, "/var/www/html/received_data.json")
    if err != nil {
      fmt.Println("Error updating JSON data in file:", err)
      continue
    }

    fmt.Println("Data updated in /var/www/html/received_data.json")
  }
}

// Fonction pour mettre à jour les données JSON dans le fichier
func updateJSONInFile(newData []byte, filePath string) error {
  // Lire les données JSON existantes
  existingData, err := ioutil.ReadFile(filePath)
  if err != nil && !os.IsNotExist(err) {
    return err
  }

  // Ajouter un saut de ligne si le fichier n'est pas vide et n'existe pas
  if len(existingData) > 0 {
    newData = append([]byte("\n"), newData...)
  }

  // Ajouter les nouvelles données JSON à la fin des données existantes
  updatedData := append(existingData, newData...)

  // Réécrire le fichier avec les données mises à jour
  err = ioutil.WriteFile(filePath, updatedData, 0644)
  if err != nil {
    return err
  }

  return nil
}

// Fonction principale
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

    go handleConnection(conn) // Gérer la connexion dans une goroutine
  }
}

// Fonction pour obtenir l'adresse IP sortante
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
