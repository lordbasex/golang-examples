package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func getLinuxDistribution() (string, error) {
	data, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "VERSION_ID") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.Trim(parts[1], `"`), nil
			}
		}
	}
	return "", nil
}

func generateMachineID() (string, error) {

	info, err1 := host.Info()
	if err1 != nil {
		return "", err1
	}

	log.Print(info.Platform)
	log.Print(info.KernelArch)
	log.Print(info.VirtualizationSystem)

	// Recopilación de metadatos del sistema
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return "", err
	}
	cpuName := cpuInfo[0].ModelName
	numCPUs := len(cpuInfo)

	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	ramSize := virtualMemory.Total / (1024 * 1024) // Convertir a MB

	hostInfo, err := host.Info()
	if err != nil {
		return "", err
	}
	osName := hostInfo.OS
	platform := hostInfo.Platform

	linuxDistribution, err := getLinuxDistribution()
	if err != nil {
		return "", err
	}

	var metadata []string
	for _, iface := range interfaces {
		metadata = append(metadata, iface.Name, iface.HardwareAddr.String())
	}
	metadata = append(metadata, hostname, osName, platform, linuxDistribution, cpuName, fmt.Sprintf("%d", numCPUs), fmt.Sprintf("%dMB", ramSize))

	// Agregar info.KernelArch e info.VirtualizationSystem a la lista de metadatos
	metadata = append(metadata, info.KernelArch, info.VirtualizationSystem)

	// Generación de una cadena única
	uniqueString := strings.Join(metadata, "|")

	log.Print(uniqueString)

	// Algoritmo de hash (MD5 en este caso)
	hash := md5.Sum([]byte(uniqueString))

	// Transformación del hash y formato del identificador
	machineID := fmt.Sprintf("%x", hash)
	return machineID, nil
}

func main() {
	machineID, err := generateMachineID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Machine ID:", machineID)
}
