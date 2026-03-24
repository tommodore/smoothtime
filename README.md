# smoothtime

**NTP without the twice-yearly trauma.**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommodore/smoothtime)](https://goreportcard.com/report/github.com/tommodore/smoothtime)
![GitHub stars](https://img.shields.io/github/stars/tommodore/smoothtime?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/tommodore/smoothtime)

A public NTP server that replaces hard Daylight Saving Time jumps with a **smooth linear drift** over months.

No more sudden clock changes on servers, homelabs, IoT devices or anywhere precise time matters.

### Live Website
→ [https://smoothtime.io](https://smoothtime.io)

### Status Page
→ [https://status.smoothtime.io](https://status.smoothtime.io)

## Supported Regions

Choose the server closest to your actual timezone:

| Region                  | Subdomain                          | NTP Command |
|-------------------------|------------------------------------|-------------|
| UK (WET/WEST)           | `europe-west.smoothtime.io`        | `server europe-west.smoothtime.io iburst` |
| Central Europe (CET/CEST) | `europe-central.smoothtime.io`   | `server europe-central.smoothtime.io iburst` |
| Eastern Europe (EET/EEST) | `europe-east.smoothtime.io`      | `server europe-east.smoothtime.io iburst` |
| US Eastern (EST/EDT)    | `us-eastern.smoothtime.io`         | `server us-eastern.smoothtime.io iburst` |
| US Central (CST/CDT)    | `us-central.smoothtime.io`         | `server us-central.smoothtime.io iburst` |
| US Mountain (MST/MDT)   | `us-mountain.smoothtime.io`        | `server us-mountain.smoothtime.io iburst` |
| US Pacific (PST/PDT)    | `us-pacific.smoothtime.io`         | `server us-pacific.smoothtime.io iburst` |

## Quick Start

```bash
# Local test (no sudo needed)
go run . -region cet

# Real NTP port
sudo go run . -region cet -ntp-port 123
```

Test the server with:
```bash
sntp -d 127.0.0.1                  # macOS
ntpdate -q 127.0.0.1 -p 1123       # Linux
```

## Features

- Smooth linear DST drift (no hard jumps)
- Correct DST schedules for all supported regions
- Stratum 2 NTP server
- HTTP `/health` endpoint
- Docker and systemd ready
- Clean status page (UptimeFlare)

## How the smooth offset works

After summer time ends, the offset gradually increases linearly until it reaches exactly +1 hours (or the equivalent for other zones) on the next DST start date — with **no jumps**.

## Deployment

See:
- `Dockerfile` + `docker-compose.yml`
- systemd service examples

## Contributing

Pull requests are welcome! Especially:
- Additional regions (Australia, New Zealand, Chile…)
- Prometheus metrics
- Improved status page

## License

MIT © 2026 tommodore & contributors

See [LICENSE](LICENSE) for details.
```

**Ready to use** — just replace the entire content of your `README.md` with the text above.

Would you like me to also update the Go code / config to match these new subdomain names (`europe-central`, `europe-west`, etc.)?
