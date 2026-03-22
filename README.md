# smoothtime

**NTP without the twice-yearly trauma.**

A public NTP server that uses a smooth sinusoidal drift instead of hard DST jumps.

No more "the clock just jumped" moments on your servers, homelab, or IoT devices.

## Features

- Sinusoidal smooth DST transition (no jumps)
- Multiple regions (CET/CEST, EET/EEST + easy to extend)
- Stratum 2 NTP server (UDP 123 or custom port)
- HTTP /health endpoint for monitoring
- Graceful shutdown
- Docker & systemd ready
- Zero external dependencies (pure Go)

## Quick start

```bash
go run . -region cet          # test port 1123
sudo go run . -region cet -ntp-port 123   # real NTP port
