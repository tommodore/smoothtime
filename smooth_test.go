package main

import (
	"math"
	"testing"
	"time"
)

func TestCalculateSmoothOffset(t *testing.T) {
	cases := []struct {
		region, at string
		min, max   float64
	}{
		{"europe-central", "2026-10-31T14:00:00Z", 1.0, 1.1},     // just after 2026-10-25 DST end
		{"europe-central", "2027-03-27T12:00:00Z", 1.95, 2.0},    // just before DST start
		{"australia-sydney", "2026-10-31T14:00:00Z", 10.8, 11.0}, // early summer
		{"australia-sydney", "2027-02-15T00:00:00Z", 10.0, 11.0}, // summer spanning New Year
		{"new-zealand", "2027-02-15T00:00:00Z", 12.0, 13.0},
		{"argentina", "2027-02-15T00:00:00Z", -3.0, -3.0},
		{"argentina", "2026-07-01T00:00:00Z", -3.0, -3.0},
	}
	for _, c := range cases {
		at, _ := time.Parse(time.RFC3339, c.at)
		got := CalculateSmoothOffset(at, c.region)
		if got < c.min-1e-9 || got > c.max+1e-9 || math.IsNaN(got) {
			t.Errorf("%s @ %s: got %.4f, want [%.2f, %.2f]", c.region, c.at, got, c.min, c.max)
		}
	}
}
