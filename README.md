# smoothtime

The first public NTP server that uses a smooth drift instead of hard +1 h / -1 h jumps.

It replaces hard Daylight Saving Time jumps with a smooth linear drift over months to prevent sudden clock changes on servers, homelabs, and IoT devices.

**Live Website** → https://smoothtime.io

## Supported Regions (11 total)

### Europe
- `europe-west.smoothtime.io` — Western Europe (WET/WEST)
- `europe-central.smoothtime.io` — Central Europe (CET/CEST)
- `europe-east.smoothtime.io` — Eastern Europe (EET/EEST)

### United States
- `us-eastern.smoothtime.io` — US Eastern (EST/EDT)
- `us-central.smoothtime.io` — US Central (CST/CDT)
- `us-mountain.smoothtime.io` — US Mountain (MST/MDT)
- `us-pacific.smoothtime.io` — US Pacific (PST/PDT)

### Southern Hemisphere (New!)
- `australia-sydney.smoothtime.io` — Australia Sydney (AEST/AEDT)
- `new-zealand.smoothtime.io` — New Zealand (NZST/NZDT)
- `chile.smoothtime.io` — Chile (CLT/CLST)
- `argentina.smoothtime.io` — Argentina (ART — constant offset, no DST)

## Quick Start

### Local Testing
```bash
go run . -region cet
```

### Real NTP Server
```bash
sudo go run . -region cet -ntp-port 123
```

## How the smooth offset works

After summer time ends, the offset increases linearly until it reaches exactly +2 hours (or the equivalent for other zones) on the next DST start date, ensuring no jumps.

The formula used is:
`offset = base + (days_passed / total_days) * 1.0`

- `base`: Winter offset (e.g. 1.0 for CET)
- `days_passed`: Days since the last DST end
- `total_days`: Number of days between the last DST end and the next DST start

## Features
- Smooth linear DST drift (no hard jumps)
- Correct DST schedules for all supported regions (including Southern Hemisphere)
- Stratum 2 NTP server
- HTTP `/health` endpoint
- Docker and systemd ready
- Clean status page (UptimeFlare)

## Deployment
- Docker: `Dockerfile` + `docker-compose.yml`
- Systemd: examples in repo
- Multi-region: use `-binds` flag (e.g. `::=europe-central,::1=australia-sydney`)

## Contributing
PRs welcome! Especially for:
- Additional regions (more Australian zones, South Africa, Brazil, etc.)
- Prometheus metrics
- Improved status page
- Better Southern Hemisphere DST edge-case handling

## License
MIT © 2026 tommodore & contributors