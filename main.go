package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// --- Configuration & Constants ---

type Rule string

const (
	RuleEU Rule = "eu"
	RuleUS Rule = "us"
)

type RegionConfig struct {
	Name string
	Base float64
	Rule Rule
}

// Regions defines the timezone offsets and DST rules for each supported node.
var Regions = map[string]RegionConfig{
	"europe-west":    {"Europe Western (WET/WEST)", 0.0, RuleEU},
	"europe-central": {"Europe Central (CET/CEST)", 1.0, RuleEU},
	"europe-east":    {"Europe Eastern (EET/EEST)", 2.0, RuleEU},
	"us-eastern":     {"US Eastern (EST/EDT)", -5.0, RuleUS},
	"us-central":     {"US Central (CST/CDT)", -6.0, RuleUS},
	"us-mountain":    {"US Mountain (MST/MDT)", -7.0, RuleUS},
	"us-pacific":     {"US Pacific (PST/PDT)", -8.0, RuleUS},
}

const ntpEpochOffset = 2208988800

// --- Main Entry Point ---

func main() {
	// Flags for binding specific IPs to regions and setting the NTP port.
	bindsFlag := flag.String("binds", "::=europe-central", "Comma-separated list of IP=region bindings")
	portFlag := flag.Int("ntp-port", 123, "UDP port to listen on for NTP")
	healthPort := flag.String("health-port", "8080", "TCP port for HTTP health checks")
	flag.Parse()

	// 1. Start the HTTP Health Check Server (for Cloudflare Uptime)
	go startHealthServer(*healthPort)

	// 2. Parse and start the regional NTP listeners
	binds := strings.Split(*bindsFlag, ",")
	if len(binds) == 0 {
		log.Fatal("No bindings specified via -binds flag")
	}

	for _, b := range binds {
		parts := strings.Split(b, "=")
		if len(parts) != 2 {
			log.Fatalf("Invalid bind format: %s. Expected IP=region", b)
		}
		ip, region := parts[0], parts[1]

		if _, exists := Regions[region]; !exists {
			log.Fatalf("Unknown region specified in bind: %s", region)
		}

		// Each IP/Region pair gets its own isolated UDP listener routine.
		go startServer(ip, *portFlag, region)
	}

	// Keep the main thread alive indefinitely.
	select {}
}

// --- Server Components ---

func startHealthServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Printf("Health Check: HTTP server active on port %s at /health", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start health server: %v", err)
	}
}

func startServer(ip string, port int, region string) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s:%d - %v", ip, port, err)
	}
	defer conn.Close()

	log.Printf("Node Online: [%s] listening on %s:%d", region, ip, port)

	buf := make([]byte, 48)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 48 {
			continue
		}
		// Handle each incoming NTP request in a non-blocking goroutine.
		go handleNTPRequest(conn, clientAddr, buf, region)
	}
}

func handleNTPRequest(conn *net.UDPConn, clientAddr *net.UDPAddr, req []byte, region string) {
	resp := make([]byte, 48)
	resp[0] = 0x1C // LI=0, VN=3, Mode=4 (Server)
	resp[1] = 2    // Stratum 2
	resp[2] = req[2]
	resp[3] = 0xFA // Precision (-6)

	copy(resp[24:32], req[40:48]) // Copy Transmit Timestamp to Originate Timestamp

	nowUTC := time.Now().UTC()
	smoothTime := ApplySmoothTime(nowUTC, region)

	sec := uint32(smoothTime.Unix() + ntpEpochOffset)
	frac := uint32((smoothTime.Nanosecond() * int(1<<32)) / 1e9)

	binary.BigEndian.PutUint32(resp[32:36], sec) // Receive Timestamp
	binary.BigEndian.PutUint32(resp[36:40], frac)
	binary.BigEndian.PutUint32(resp[40:44], sec) // Transmit Timestamp
	binary.BigEndian.PutUint32(resp[44:48], frac)

	conn.WriteToUDP(resp, clientAddr)
}

// --- Time Manipulation Logic ---

func getNthSunday(year int, month time.Month, nth int) time.Time {
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	count := 0
	for {
		if t.Weekday() == time.Sunday {
			count++
			if count == nth {
				return t
			}
		}
		t = t.AddDate(0, 0, 1)
	}
}

func getLastSunday(year int, month time.Month) time.Time {
	t := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func getDstDates(year int, rule Rule) (start time.Time, end time.Time) {
	if rule == RuleEU {
		return getLastSunday(year, time.March), getLastSunday(year, time.October)
	} else if rule == RuleUS {
		return getNthSunday(year, time.March, 2), getNthSunday(year, time.November, 1)
	}
	return time.Time{}, time.Time{}
}

func CalculateSmoothOffset(now time.Time, regionID string) float64 {
	cfg := Regions[regionID]
	year := now.Year()
	tStart, tEnd := getDstDates(year, cfg.Rule)

	var offset float64
	if !now.Before(tStart) && now.Before(tEnd) {
		totalSummer := tEnd.Sub(tStart).Seconds()
		passedSummer := now.Sub(tStart).Seconds()
		fraction := passedSummer / totalSummer
		offset = cfg.Base + 1.0 - fraction
	} else {
		wStart := tEnd
		wEnd, _ := getDstDates(year+1, cfg.Rule)
		if now.Before(tStart) {
			_, prevEnd := getDstDates(year-1, cfg.Rule)
			wStart = prevEnd
			wEnd = tStart
		}
		totalWinter := wEnd.Sub(wStart).Seconds()
		passedWinter := now.Sub(wStart).Seconds()
		fraction := passedWinter / totalWinter
		offset = cfg.Base + fraction
	}
	return offset
}

func ApplySmoothTime(utcNow time.Time, regionID string) time.Time {
	offsetHours := CalculateSmoothOffset(utcNow, regionID)
	offsetDuration := time.Duration(offsetHours * float64(time.Hour))
	return utcNow.Add(offsetDuration)
}
