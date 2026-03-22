package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"math"
)

const defaultNTPPort = "1123"
const defaultHTTPPort = "8080"

type RegionConfig struct {
	Name      string
	Base      float64
	Amplitude float64
	Phase     int
}

var regions = map[string]RegionConfig{
	"cet":  {"Central Europe (CET/CEST)", 1.0, 0.5, 80},
	"eest": {"Eastern Europe (EET/EEST)", 2.0, 0.5, 80},
}

func smoothOffset(region string) float64 {
	cfg, ok := regions[region]
	if !ok {
		cfg = regions["cet"]
		slog.Warn("unknown region, falling back to CET", "region", region)
	}
	doy := time.Now().UTC().YearDay()
	phase := 2 * math.Pi * float64(doy-cfg.Phase) / 365.0
	return cfg.Base + cfg.Amplitude*(1+math.Sin(phase))
}

func toNTPTime(t time.Time) [8]byte {
	sec := uint32(t.Unix() + 2208988800)
	frac := uint32(t.Nanosecond() * 4294967296 / 1e9)
	var b [8]byte
	binary.BigEndian.PutUint32(b[0:4], sec)
	binary.BigEndian.PutUint32(b[4:8], frac)
	return b
}

func main() {
	regionPtr := flag.String("region", "cet", "region (cet, eest)")
	ntpPortPtr := flag.String("ntp-port", defaultNTPPort, "NTP UDP port")
	httpPortPtr := flag.String("http-port", defaultHTTPPort, "HTTP health port")
	flag.Parse()

	region := *regionPtr
	offsetHours := smoothOffset(region)

	slog.Info("smoothtime started",
		"region", regions[region].Name,
		"offset_hours", fmt.Sprintf("%.3f", offsetHours),
		"ntp_port", *ntpPortPtr,
		"http_port", *httpPortPtr)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// HTTP health check (optional but very useful on VPS)
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"ok","region":"%s","offset":%.3f}`, region, offsetHours)
		})
		if err := http.ListenAndServe(":"+*httpPortPtr, nil); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
		}
	}()

	// NTP server (exactly the same reliable logic you already tested)
	addr, _ := net.ResolveUDPAddr("udp", ":"+*ntpPortPtr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 48)
	for {
		select {
		case <-ctx.Done():
			slog.Info("shutting down gracefully")
			return
		default:
			n, client, err := conn.ReadFromUDP(buf)
			if err != nil || n < 48 || (buf[0]&0x07) != 3 {
				continue
			}

			now := time.Now().UTC()
			smoothTime := now.Add(time.Duration(offsetHours*3600) * time.Second)

			resp := make([]byte, 48)
			resp[0] = 0x24
			resp[1] = 2
			resp[2] = 6
			resp[3] = 0xec

			ref := toNTPTime(smoothTime.Add(-time.Second))
			copy(resp[16:24], ref[:])
			copy(resp[24:32], buf[40:48])
			rxtx := toNTPTime(smoothTime)
			copy(resp[32:40], rxtx[:])
			copy(resp[40:48], rxtx[:])

			conn.WriteToUDP(resp, client)
		}
	}
}
