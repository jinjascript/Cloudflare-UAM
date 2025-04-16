// Jinja CPU Stress Test Utility
// Part of the Jinja CPU Monitoring System
// Created by Jinja (2025)
// This utility generates high CPU load for testing the monitoring system

package main

import (
        "fmt"
        "runtime"
        "time"
)

func main() {
        // Use all available CPU cores
        numCPUs := runtime.NumCPU()
        fmt.Printf("Starting CPU stress test on %d cores\n", numCPUs)

        // Run for 60 seconds
        duration := 60 * time.Second
        endTime := time.Now().Add(duration)

        // Start CPU-intensive tasks on each core
        for i := 0; i < numCPUs; i++ {
                go func(id int) {
                        fmt.Printf("Starting worker %d\n", id)
                        // Perform CPU-intensive calculations until time is up
                        for time.Now().Before(endTime) {
                                // Perform useless but CPU-intensive calculations
                                for j := 0; j < 10000000; j++ {
                                        _ = j * j / (j + 1)
                                }
                        }
                }(i)
        }

        // Wait until the end time
        fmt.Println("Stress test running, press Ctrl+C to stop early...")
        time.Sleep(duration)
        fmt.Println("Stress test completed")
}