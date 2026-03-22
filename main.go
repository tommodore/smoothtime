package main

import (
	"encoding/binary"
	"flag"
	"log"
	"net"
	"strings"
	"time"
)

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

func main() {
	// The new binds flag allows mapping specific IPs to specific mathematical regions
	bindsFlag := flag.String("binds", "0.0.0.0=europe-central", "Comma-separated list of IP=region bindings (e.g. 10.0.0.1=europe-central,10.0.0.2=us-eastern)")
	portFlag := flag.Int("ntp-port", 123, "UDP port to listen on")
	flag.Parse()

	binds := strings.Split(*bindsFlag, ",")
	if len(binds) == 0 {
		log.Fatal("No bindings specified")
	}

	// Spin up a concurrent UDP listener for every IP provided
	for _, b := range binds {
		parts := strings.Split(b, "=")
		if len(parts) != 2 {
			log.Fatalf("Invalid bind format: %s. Expected IP=region", b)
		}
		ip, region := parts[0], parts[1]

		if _, exists := Regions[region]; !exists {
			log.Fatalf("Unknown region specified in bind: %s", region)
		}

		go startServer(ip, *portFlag, region)
	}

	// Block the main thread forever while the goroutines handle the UDP traffic
	select {}
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

	log.Printf("Node Active: Listening on %s:%d routing to [%s]", ip, port, region)

	buf := make([]byte, 48)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n < 48 {
			continue
		}
		// Pass the specific region of this listener into the request handler
		go handleNTPRequest(conn, clientAddr, buf, region)
	}
}

func handleNTPRequest(conn *net.UDPConn, clientAddr *net.UDPAddr, req []byte, region string) {
	resp := make([]byte, 48)
	
	resp[0] = 0x1C   // LI=0, VN=3, Mode=4 (Server)
	resp[1] = 2      // Stratum 2
	resp[2] = req[2] // Poll interval
	resp[3] = 0xFA   // Precision (-6)

	copy(resp[24:32], req[40:48])

	nowUTC := time.Now().UTC()
	smoothTime := ApplySmoothTime(nowUTC, region)

	sec := uint32(smoothTime.Unix() + ntpEpochOffset)
	frac := uint32((smoothTime.Nanosecond() * int(1<<32)) / 1e9)

	binary.BigEndian.PutUint32(resp[32:36], sec)
	binary.BigEndian.PutUint32(resp[36:40], frac)
	binary.BigEndian.PutUint32(resp[40:44], sec)
	binary.BigEndian.PutUint32(resp[44:48], frac)

	conn.WriteToUDP(resp, clientAddr)
}

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
