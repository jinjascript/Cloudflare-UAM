# Jinja CPU Monitoring System
> Cloudflare UAM Script

	
***This program is a CPU monitoring system developed using the `Go` programming language. The system constantly monitors your server's CPU usage and _automatically_ activates Cloudflare's "Under Attack Mode" feature when high CPU usage is detected, protecting against DDoS attacks.***

## Working Principle
***The program checks the CPU usage every 10 seconds.
When the CPU usage exceeds 80% and 3 consecutive high readings are detected, Cloudflare's "Under Attack Mode" feature is activated.
When the CPU usage returns to normal and 5 consecutive normal readings are detected, Cloudflare protection mode is disabled.
All processes and CPU measurements are recorded via the console.
This system is ideal for automatic protection, especially on high-traffic websites or servers at risk of DDoS attacks.***

Run the program:
```
go run main.go
```
Configuration Settings (Optional):
```rust
CF_API_KEY=your_actual_cloudflare_api_key
CF_EMAIL=your_actual_cloudflare_email
CF_ZONE_ID=your_actual_cloudflare_zone_id
}
```
