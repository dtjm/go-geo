package geo

import (
	"math"
)

const (
	EarthRadiusMi = 3959
	EarthRadiusKm = 6371
)

// Calculate Haversine distance between two lat/lon points on Earth.
// Takes parameters in radians
// Returns result in meters
func Haversine(startLat, startLon, endLat, endLon float64) float64 {
	dLat := endLat - startLat
	dLon := endLon - startLon

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(startLat)*math.Cos(endLat)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := EarthRadiusKm * c
	return d
}
