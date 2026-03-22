# smoothtime

**NTP without the twice-yearly trauma.**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommodore/smoothtime)](https://goreportcard.com/report/github.com/tommodore/smoothtime)
![GitHub stars](https://img.shields.io/github/stars/tommodore/smoothtime?style=social)
![GitHub forks](https://img.shields.io/github/forks/tommodore/smoothtime?style=social)
![GitHub issues](https://img.shields.io/github/issues/tommodore/smoothtime)
![GitHub last commit](https://img.shields.io/github/last-commit/tommodore/smoothtime)

A public NTP server that replaces the traditional hard +1 h / -1 h Daylight Saving Time jumps with a **smooth sinusoidal drift** spread over months.

No more sudden clock jumps on your servers, homelab, IoT devices, or anywhere NTP is used.

## Features

- Smooth sinusoidal DST transition (no discontinuities)
- Supported regions: Central Europe (CET/CEST), Eastern Europe (EET/EEST) — easy to extend
- Stratum 2 NTP server (UDP 123 or custom port for testing)
- Optional HTTP `/health` endpoint for monitoring
- Graceful shutdown (SIGINT / SIGTERM)
- Pure Go — zero external dependencies
- Docker-ready
- systemd service example included

## How the smooth offset works

The offset is calculated using this formula:

```math
offset = base + amplitude \times (1 + \sin(2\pi \times (doy - phase) / 365))
```

- `base` — winter UTC offset (e.g. 1.0 for CET)
- `amplitude` — half the DST swing (usually 0.5 → ±1 hour)
- `phase` — shifts the peak to align with real summer time
- `doy` — day of year (1–365/366)

This produces a gentle, continuous curve instead of abrupt steps.

## Quick Start (local testing)

```bash
# Test on non-privileged port 1123 (no sudo needed)
go run . -region cet

# Eastern Europe region
go run . -region eest

# Real NTP port (requires sudo)
sudo go run . -region cet -ntp-port 123
```

Test with:

```bash
# macOS
sntp -d 127.0.0.1

# Linux
ntpdate -q 127.0.0.1 -p 1123
chronyd -q 'server 127.0.0.1 port 1123 iburst'
```

## Installation & Deployment

### From source

```bash
git clone https://github.com/tommodore/smoothtime.git
cd smoothtime
go build -o smoothtime
sudo ./smoothtime -region cet -ntp-port 123
```

### Docker

```bash
docker build -t tommodore/smoothtime .
docker run --network host --cap-add=NET_BIND_SERVICE tommodore/smoothtime -region cet
```

See `docker-compose.yml` for multi-region setup.

### systemd (production)

Create `/etc/systemd/system/smoothtime-cet.service`:

```ini
[Unit]
Description=smoothtime NTP Server – Central Europe
After=network.target

[Service]
ExecStart=/usr/local/bin/smoothtime -region cet -ntp-port 123
Restart=always
User=nobody
Group=nogroup
LimitNOFILE=65535
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now smoothtime-cet
```

## Monitoring

HTTP health check:  
`http://your-server:8080/health`

## Contributing

Pull requests welcome! Especially:
- More regions (US, UK, Australia, New Zealand, Chile…)
- Prometheus metrics
- Tests
- Landing page / website

## License

MIT © 2026 tommodore & contributors

See [LICENSE](LICENSE) for details.
```

Copy everything above (including the first `# smoothtime` line) and paste it directly into your `README.md`.  
The badges will work instantly after you push.  

You’re all set! 🚀  

Ready for the next step (CI workflow, more regions, or launch plan)? Just say the word.
