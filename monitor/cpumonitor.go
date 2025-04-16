// Jinja CPU Monitor
// Part of the Jinja CPU Monitoring System
// Created by Jinja (2025)
// This module provides CPU usage monitoring functionality

package monitor

import (
        "fmt"
        "time"

        "github.com/shirou/gopsutil/v3/cpu"
)

// GetCPULoad returns the current CPU load as a percentage (0-100)
func GetCPULoad() (float64, error) {
        // Get CPU usage percentage for 1 second
        percentages, err := cpu.Percent(time.Second, false)
        if err != nil {
                return 0, fmt.Errorf("failed to get CPU utilization: %w", err)
        }

        if len(percentages) == 0 {
                return 0, fmt.Errorf("no CPU utilization data available")
        }

        // Return the overall CPU usage
        return percentages[0], nil
}

// GetDetailedCPUUsage returns detailed CPU statistics
func GetDetailedCPUUsage() (map[string]float64, error) {
        stats, err := cpu.Times(false)
        if err != nil {
                return nil, fmt.Errorf("failed to get detailed CPU stats: %w", err)
        }

        if len(stats) == 0 {
                return nil, fmt.Errorf("no CPU statistics available")
        }

        // Extract the first CPU's stats
        stat := stats[0]
        
        total := stat.User + stat.System + stat.Idle + stat.Nice + 
                stat.Iowait + stat.Irq + stat.Softirq + stat.Steal + stat.Guest + stat.GuestNice

        if total <= 0 {
                return nil, fmt.Errorf("invalid CPU statistics (total <= 0)")
        }

        // Calculate percentages
        return map[string]float64{
                "user":       (stat.User / total) * 100,
                "system":     (stat.System / total) * 100,
                "idle":       (stat.Idle / total) * 100,
                "nice":       (stat.Nice / total) * 100,
                "iowait":     (stat.Iowait / total) * 100,
                "irq":        (stat.Irq / total) * 100,
                "softirq":    (stat.Softirq / total) * 100,
                "steal":      (stat.Steal / total) * 100,
                "guest":      (stat.Guest / total) * 100,
                "guestNice":  (stat.GuestNice / total) * 100,
        }, nil
}
