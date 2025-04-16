// Jinja CPU Monitoring System
// Developed by Jinja (2025)
// A lightweight CPU monitoring system that integrates with Cloudflare to protect against DDoS attacks
// Version: 1.0

package main

import (
        "log"
        "os"
        "os/signal"
        "syscall"
        "time"

        "github.com/joho/godotenv"

        "ddos-protection/cloudflare"
        "ddos-protection/config"
        "ddos-protection/monitor"
)

func main() {
        // Load environment variables from .env file
        err := godotenv.Load()
        if err != nil {
                log.Println("Warning: .env file not found, will use system environment variables")
        }

        // Initialize configuration
        cfg, err := config.LoadConfig()
        if err != nil {
                log.Fatalf("Failed to load configuration: %v", err)
        }

        // Create a cloudflare client
        cfClient := cloudflare.NewClient(cfg.CloudflareAPIKey, cfg.CloudflareEmail, cfg.CloudflareZoneID)

        // Create a channel to handle shutdown signals
        sigs := make(chan os.Signal, 1)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

        // Create a done channel to signal when the program is finished
        done := make(chan bool, 1)

        // Set initial attack mode state
        attackModeEnabled := false

        // Start monitoring in a goroutine
        go func() {
                // Create ticker for regular CPU checks
                ticker := time.NewTicker(time.Duration(cfg.MonitoringInterval) * time.Second)
                defer ticker.Stop()

                log.Printf("Starting CPU monitoring (threshold: %d%%, interval: %ds)\n", 
                        cfg.CPUThreshold, cfg.MonitoringInterval)

                consecutiveHighLoads := 0
                consecutiveLowLoads := 0

                for {
                        select {
                        case <-ticker.C:
                                cpuLoad, err := monitor.GetCPULoad()
                                if err != nil {
                                        log.Printf("Error getting CPU load: %v", err)
                                        continue
                                }

                                log.Printf("Current CPU load: %.2f%%", cpuLoad)

                                // Check if CPU load exceeds threshold
                                if cpuLoad >= float64(cfg.CPUThreshold) {
                                        consecutiveHighLoads++
                                        consecutiveLowLoads = 0
                                        log.Printf("High CPU load detected (%d consecutive readings)", consecutiveHighLoads)

                                        // Enable "Under Attack Mode" if not already enabled and we've had enough consecutive high readings
                                        if !attackModeEnabled && consecutiveHighLoads >= cfg.RequiredHighReadings {
                                                log.Println("Suspected DDoS attack! Enabling Cloudflare 'Under Attack Mode'...")
                                                err := cfClient.EnableUnderAttackMode()
                                                if err != nil {
                                                        log.Printf("Error enabling 'Under Attack Mode': %v", err)
                                                } else {
                                                        attackModeEnabled = true
                                                        log.Println("Cloudflare 'Under Attack Mode' enabled successfully")
                                                }
                                        }
                                } else {
                                        consecutiveLowLoads++
                                        consecutiveHighLoads = 0
                                        
                                        // Disable "Under Attack Mode" if enabled and we've had enough consecutive low readings
                                        if attackModeEnabled && consecutiveLowLoads >= cfg.RequiredLowReadings {
                                                log.Println("CPU load has returned to normal. Disabling 'Under Attack Mode'...")
                                                err := cfClient.DisableUnderAttackMode()
                                                if err != nil {
                                                        log.Printf("Error disabling 'Under Attack Mode': %v", err)
                                                } else {
                                                        attackModeEnabled = false
                                                        log.Println("Cloudflare 'Under Attack Mode' disabled successfully")
                                                }
                                        }
                                }
                        case <-done:
                                return
                        }
                }
        }()

        // Wait for termination signal
        <-sigs
        log.Println("Shutting down...")

        // If attack mode is enabled, disable it before exiting
        if attackModeEnabled {
                log.Println("Attempting to disable 'Under Attack Mode' before exiting...")
                err := cfClient.DisableUnderAttackMode()
                if err != nil {
                        log.Printf("Error disabling 'Under Attack Mode': %v", err)
                } else {
                        log.Println("Cloudflare 'Under Attack Mode' disabled successfully")
                }
        }

        // Signal the monitoring routine to stop
        done <- true
        log.Println("CPU monitoring stopped. Goodbye!")
}
