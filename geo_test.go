package geo

import (
	"math"
	"testing"
)

const (
	DistThreshold         = 0.01
	DistThresholdVincenty = 0.5
)

func Deg2Rad(x float64) float64 {
	return x * math.Pi / 180
}

func TestHaversine(t *testing.T) {
	locations := make(map[string]([2]float64))
	locations["Google HQ"] = [2]float64{37.422045, -122.084347}
	locations["San Francisco"] = [2]float64{37.77493, -122.419416}
	locations["Eiffel Tower"] = [2]float64{48.8582, 2.294407}
	locations["Sydney Opera House"] = [2]float64{-33.856553, 151.214696}

	testCases := make(map[string](map[string]float64))
	for loc := range locations {
		testCases[loc] = make(map[string]float64)
	}

	testCases["Google HQ"]["San Francisco"] = 49.103
	testCases["Google HQ"]["Eiffel Tower"] = 8967.042
	testCases["Google HQ"]["Sydney Opera House"] = 11952.717

	for startLoc, cases := range testCases {
		start := locations[startLoc]
		for endLoc, expected := range cases {
			end := locations[endLoc]
			dist := Haversine(Deg2Rad(start[0]), Deg2Rad(start[1]),
				Deg2Rad(end[0]), Deg2Rad(end[1]))
			if dist < expected-DistThreshold || dist > expected+DistThreshold {
				t.Errorf("Distance from %s to %s should be ~%v km, got %v km",
					startLoc, endLoc, expected, dist)
			}
		}
	}
}

func BenchmarkHaversine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Haversine(1.1, 2.2, 3.3, 4.4)
	}
}

func TestDistVincenty(t *testing.T) {
	locations := make(map[string]([2]float64))
	locations["Google HQ"] = [2]float64{37.422045, -122.084347}
	locations["San Francisco"] = [2]float64{37.77493, -122.419416}
	locations["Eiffel Tower"] = [2]float64{48.8582, 2.294407}
	locations["Sydney Opera House"] = [2]float64{-33.856553, 151.214696}

	testCases := make(map[string](map[string]float64))
	for loc := range locations {
		testCases[loc] = make(map[string]float64)
	}

	// http://www.wolframalpha.com/input/?i=distance+between+37.422045N,+122.084347W+and+37.77493N,+122.419416W
	testCases["Google HQ"]["San Francisco"] = 49.09
	// http://www.wolframalpha.com/input/?i=distance+between+37.422045N,+122.084347W+and+48.8582N,+2.294E
	testCases["Google HQ"]["Eiffel Tower"] = 8990.0
	// http://www.wolframalpha.com/input/?i=distance+between+37.422045N,+122.084347W+and+33.856553S,+151.214696E
	testCases["Google HQ"]["Sydney Opera House"] = 11940.0

	for startLoc, cases := range testCases {
		start := locations[startLoc]
		for endLoc, expected := range cases {
			end := locations[endLoc]
			dist := DistVincenty(Deg2Rad(start[0]), Deg2Rad(start[1]),
				Deg2Rad(end[0]), Deg2Rad(end[1]))
			if dist < expected-DistThresholdVincenty || dist > expected+DistThresholdVincenty {
				t.Errorf("Distance from %s to %s should be ~%v km, got %v km",
					startLoc, endLoc, expected, dist)
			}
		}
	}
}

func BenchmarkDistVincenty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistVincenty(1.1, 2.2, 3.3, 4.4)
	}
}
