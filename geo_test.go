package geo

import (
	"math"
	"testing"
)

func Deg2Rad(x float64) float64 {
	return x * math.Pi / 180
}

func DistTest(algorithm string, t *testing.T) {

	locations := make(map[string]([2]float64))
	locations["Google HQ"] = [2]float64{37.422045, -122.084347}
	locations["San Francisco"] = [2]float64{37.77493, -122.419416}
	locations["Eiffel Tower"] = [2]float64{48.8582, 2.294407}
	locations["Sydney Opera House"] = [2]float64{-33.856553, 151.214696}

	testCases := make(map[string](map[string]float64))
	for loc := range locations {
		testCases[loc] = make(map[string]float64)
	}

	switch algorithm {
	case "Haversine":
		testCases["Google HQ"]["San Francisco"] = 49.103
		testCases["Google HQ"]["Eiffel Tower"] = 8967.042
		testCases["Google HQ"]["Sydney Opera House"] = 11952.717

	case "Vincenty":
		// expected distances from http://www.wolframalpha.com
		testCases["Google HQ"]["San Francisco"] = 49.09
		testCases["Google HQ"]["Eiffel Tower"] = 8990.0
		testCases["Google HQ"]["Sydney Opera House"] = 11940.0
	}

	for startLoc, cases := range testCases {
		start := locations[startLoc]
		for endLoc, expected := range cases {
			end := locations[endLoc]
			if algorithm == "Haversine" {
				const (
					DistThreshold = 0.01
				)
				dist := Haversine(Deg2Rad(start[0]), Deg2Rad(start[1]),
					Deg2Rad(end[0]), Deg2Rad(end[1]))
				if dist < expected-DistThreshold || dist > expected+DistThreshold {
					t.Errorf("Distance from %s to %s should be ~%v km, got %v km",
						startLoc, endLoc, expected, dist)
				}
			} else if algorithm == "Vincenty" {
				const (
					DistThreshold = 0.5
				)
				dist := DistVincenty(Deg2Rad(start[0]), Deg2Rad(start[1]),
					Deg2Rad(end[0]), Deg2Rad(end[1]))
				if dist < expected-DistThreshold || dist > expected+DistThreshold {
					t.Errorf("Distance from %s to %s should be ~%v km, got %v km",
						startLoc, endLoc, expected, dist)
				}
			}
		}
	}

}

func TestHaversine(t *testing.T) {
	DistTest("Haversine", t)
}

func BenchmarkHaversine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Haversine(1.1, 2.2, 3.3, 4.4)
	}
}

func TestDistVincenty(t *testing.T) {
	DistTest("Vincenty", t)
}

func BenchmarkDistVincenty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DistVincenty(1.1, 2.2, 3.3, 4.4)
	}
}
