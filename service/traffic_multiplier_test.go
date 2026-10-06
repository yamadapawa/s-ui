package service

import (
	"math"
	"testing"
)

func TestBilledTraffic(t *testing.T) {
	tests := []struct {
		name          string
		actual        int64
		multiplier    float64
		want          int64
		wantRemainder float64
	}{
		{name: "double", actual: 100, multiplier: 2, want: 200},
		{name: "fractional", actual: 101, multiplier: 1.5, want: 151, wantRemainder: 0.5},
		{name: "default for zero value", actual: 42, multiplier: 0, want: 42},
		{name: "invalid multiplier defaults to one", actual: 42, multiplier: math.NaN(), want: 42},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, remainder := billedTraffic(test.actual, test.multiplier, 0)
			if got != test.want {
				t.Fatalf("billedTraffic(%d, %v) = %d, want %d", test.actual, test.multiplier, got, test.want)
			}
			if math.Abs(remainder-test.wantRemainder) > 1e-9 {
				t.Fatalf("billedTraffic remainder = %v, want %v", remainder, test.wantRemainder)
			}
		})
	}

	remainder := 0.0
	var total int64
	for range 10 {
		charged, nextRemainder := billedTraffic(1, 0.1, remainder)
		total += charged
		remainder = nextRemainder
	}
	if total != 1 {
		t.Fatalf("ten 1-byte intervals at 0.1x billed %d bytes, want 1", total)
	}
	if remainder < 0 {
		t.Fatalf("quota remainder became negative: %v", remainder)
	}
}
