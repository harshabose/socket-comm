package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	fmt.Println("=== CPU Usage ===")
	CPU()

	fmt.Println("\n=== Memory Usage ===")
	Memory()

	fmt.Println("\n=== Disk Usage ===")
	Disk()
}

func CPU() {
	// Get CPU usage percentage for all cores combined (false parameter)
	// Wait 100ms for measurement interval
	percentages, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil {
		log.Fatalf("Error getting CPU percentages: %v", err)
	}

	// Print the result
	// Since we passed false as the second parameter, we get a single value
	// representing the average usage across all CPU cores
	fmt.Printf("CPU Usage (all cores): %.2f%%\n", percentages[0])

	// Now let's get per-core CPU usage
	perCorePercentages, err := cpu.Percent(100*time.Millisecond, true)
	if err != nil {
		log.Fatalf("Error getting per-core CPU percentages: %v", err)
	}

	// Print per-core results
	for i, percentage := range perCorePercentages {
		fmt.Printf("CPU Core #%d Usage: %.2f%%\n", i, percentage)
	}

	// Demonstrate continuous monitoring
	fmt.Println("\nContinuous CPU monitoring (5 seconds):")
	for i := 0; i < 5; i++ {
		percentages, err := cpu.Percent(1000*time.Millisecond, false)
		if err != nil {
			log.Printf("Error getting CPU percentages: %v", err)
			continue
		}
		fmt.Printf("CPU Usage at %s: %.2f%%\n", time.Now().Format("15:04:05"), percentages[0])
	}
}

func Memory() {
	// Get virtual memory statistics
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error getting virtual memory statistics: %v", err)
	}

	// Print memory usage information
	fmt.Printf("Total memory: %.2f GB\n", float64(vmStat.Total)/(1024*1024*1024))
	fmt.Printf("Available memory: %.2f GB\n", float64(vmStat.Available)/(1024*1024*1024))
	fmt.Printf("Used memory: %.2f GB\n", float64(vmStat.Used)/(1024*1024*1024))
	fmt.Printf("Memory usage percentage: %.2f%%\n", vmStat.UsedPercent)

	// Get swap memory statistics
	swapStat, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Error getting swap memory statistics: %v", err)
	} else {
		fmt.Printf("\nSwap memory total: %.2f GB\n", float64(swapStat.Total)/(1024*1024*1024))
		fmt.Printf("Swap memory used: %.2f GB\n", float64(swapStat.Used)/(1024*1024*1024))
		fmt.Printf("Swap memory usage percentage: %.2f%%\n", swapStat.UsedPercent)
	}
}

func Disk() {
	// Get disk partitions
	partitions, err := disk.Partitions(false) // false means physical partitions only
	if err != nil {
		log.Fatalf("Error getting disk partitions: %v", err)
	}

	fmt.Println("Disk partitions and usage:")
	for _, partition := range partitions {
		fmt.Printf("\nDevice: %s\n", partition.Device)
		fmt.Printf("Mount point: %s\n", partition.Mountpoint)
		fmt.Printf("File system type: %s\n", partition.Fstype)

		// Get usage statistics for this partition
		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Printf("Error getting usage statistics for %s: %v", partition.Mountpoint, err)
			continue
		}

		fmt.Printf("Total space: %.2f GB\n", float64(usageStat.Total)/(1024*1024*1024))
		fmt.Printf("Free space: %.2f GB\n", float64(usageStat.Free)/(1024*1024*1024))
		fmt.Printf("Used space: %.2f GB\n", float64(usageStat.Used)/(1024*1024*1024))
		fmt.Printf("Usage percentage: %.2f%%\n", usageStat.UsedPercent)
	}

	// Get IO counters
	ioCounters, err := disk.IOCounters()
	if err != nil {
		log.Printf("Error getting disk IO counters: %v", err)
	} else {
		fmt.Println("\nDisk IO statistics:")
		for deviceName, ioStat := range ioCounters {
			fmt.Printf("\nDevice: %s\n", deviceName)
			fmt.Printf("Read count: %d\n", ioStat.ReadCount)
			fmt.Printf("Write count: %d\n", ioStat.WriteCount)
			fmt.Printf("Read bytes: %.2f MB\n", float64(ioStat.ReadBytes)/(1024*1024))
			fmt.Printf("Write bytes: %.2f MB\n", float64(ioStat.WriteBytes)/(1024*1024))
			fmt.Printf("Read time: %d ms\n", ioStat.ReadTime)
			fmt.Printf("Write time: %d ms\n", ioStat.WriteTime)
		}
	}
}
