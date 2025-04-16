// Jinja Configuration Manager
// Part of the Jinja CPU Monitoring System
// Created by Jinja (2025)
// This module handles configuration loading from environment variables

package config

import (
        "fmt"
        "os"
        "strconv"
)

// Config stores the application configuration
type Config struct {
        // Cloudflare API credentials
        CloudflareAPIKey string
        CloudflareEmail  string
        CloudflareZoneID string
        
        // CPU monitoring settings
        CPUThreshold        int // Percentage threshold to trigger alert (default: 80%)
        MonitoringInterval  int // Time between checks in seconds (default: 10s)
        RequiredHighReadings int // Number of consecutive high readings before enabling attack mode (default: 3)
        RequiredLowReadings  int // Number of consecutive low readings before disabling attack mode (default: 5)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
        config := &Config{}

        // Load required Cloudflare API credentials
        config.CloudflareAPIKey = os.Getenv("CF_API_KEY")
        if config.CloudflareAPIKey == "" {
                return nil, fmt.Errorf("CF_API_KEY environment variable is required")
        }

        config.CloudflareEmail = os.Getenv("CF_EMAIL")
        if config.CloudflareEmail == "" {
                return nil, fmt.Errorf("CF_EMAIL environment variable is required")
        }

        config.CloudflareZoneID = os.Getenv("CF_ZONE_ID")
        if config.CloudflareZoneID == "" {
                return nil, fmt.Errorf("CF_ZONE_ID environment variable is required")
        }

        // Load CPU threshold with default value
        thresholdStr := os.Getenv("CPU_THRESHOLD")
        if thresholdStr == "" {
                config.CPUThreshold = 80 // Default to 80%
        } else {
                threshold, err := strconv.Atoi(thresholdStr)
                if err != nil {
                        return nil, fmt.Errorf("invalid CPU_THRESHOLD value: %v", err)
                }
                
                if threshold <= 0 || threshold > 100 {
                        return nil, fmt.Errorf("CPU_THRESHOLD must be between 1 and 100")
                }
                
                config.CPUThreshold = threshold
        }

        // Load monitoring interval with default value
        intervalStr := os.Getenv("MONITORING_INTERVAL")
        if intervalStr == "" {
                config.MonitoringInterval = 10 // Default to 10 seconds
        } else {
                interval, err := strconv.Atoi(intervalStr)
                if err != nil {
                        return nil, fmt.Errorf("invalid MONITORING_INTERVAL value: %v", err)
                }
                
                if interval < 1 {
                        return nil, fmt.Errorf("MONITORING_INTERVAL must be at least 1 second")
                }
                
                config.MonitoringInterval = interval
        }

        // Load required high readings count with default value
        highReadingsStr := os.Getenv("REQUIRED_HIGH_READINGS")
        if highReadingsStr == "" {
                config.RequiredHighReadings = 3 // Default to 3 readings
        } else {
                highReadings, err := strconv.Atoi(highReadingsStr)
                if err != nil {
                        return nil, fmt.Errorf("invalid REQUIRED_HIGH_READINGS value: %v", err)
                }
                
                if highReadings < 1 {
                        return nil, fmt.Errorf("REQUIRED_HIGH_READINGS must be at least 1")
                }
                
                config.RequiredHighReadings = highReadings
        }

        // Load required low readings count with default value
        lowReadingsStr := os.Getenv("REQUIRED_LOW_READINGS")
        if lowReadingsStr == "" {
                config.RequiredLowReadings = 5 // Default to 5 readings
        } else {
                lowReadings, err := strconv.Atoi(lowReadingsStr)
                if err != nil {
                        return nil, fmt.Errorf("invalid REQUIRED_LOW_READINGS value: %v", err)
                }
                
                if lowReadings < 1 {
                        return nil, fmt.Errorf("REQUIRED_LOW_READINGS must be at least 1")
                }
                
                config.RequiredLowReadings = lowReadings
        }

        return config, nil
}
