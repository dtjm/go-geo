package geo

import (
	// "errors"
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

// Rounding function courtesy David Vaini
// https://gist.github.com/DavidVaini/10308388
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// DistVicenty is a re-write of the function LatLong.distVincenty
// by Chris Veness ((c) 2002-2006 Chris Veness)
// http://www.5thandpenn.com/GeoMaps/GMapsExamples/distanceComplete2.html
// https://www.movable-type.co.uk/scripts/latlong-vincenty.html
// based on DIRECT  AND INVERSE SOLUTIONS OF GEODESICS ON THE ELLIPSOID
// WITH APPLICATION OF NESTED EQUATIONS, T. Vincenty, April 1975
// https://www.ngs.noaa.gov/PUBS_LIB/inverse.pdf
// Calculate geodesic distance (in m) between two points specified by
// latitude/longitude using Vincenty inverse formula for ellipsoids
// Returns result in metres
func DistVincenty(startLat, startLon, endLat, endLon float64) float64 {
	a := 6378137.0
	b := 6356752.3142
	f := 1 / 298.257223563 // WGS-84 ellipsiod
	L := endLon - startLon
	U1 := math.Atan((1 - f) * math.Tan(startLat))
	U2 := math.Atan((1 - f) * math.Tan(endLat))
	sinU1 := math.Sin(U1)
	cosU1 := math.Cos(U1)
	sinU2 := math.Sin(U2)
	cosU2 := math.Cos(U2)
	lambda := L
	lambdaP := 2 * math.Pi
	iterLimit := 20
	cosSqAlpha := 1.0
	sinSigma := 1.0
	cos2SigmaM := 1.0
	cosSigma := 1.0
	sigma := 1.0

	for ((math.Abs(lambda - lambdaP)) > math.Pow10(-12)) && (iterLimit > 0) {
		sinLambda := math.Sin(lambda)
		cosLambda := math.Cos(lambda)
		sinSigma = math.Sqrt((cosU2*sinLambda)*(cosU2*sinLambda) + (cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda))
		if sinSigma == 0 {
			return 0.0
		}
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha := cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - 2*sinU1*sinU2/cosSqAlpha
		if math.IsNaN(cos2SigmaM) {
			cos2SigmaM = 0.0 // equatorial line: cosSqAlpha=0 (§6)
		}
		C := f / 16 * cosSqAlpha * (4 + f*(4-3*cosSqAlpha))
		lambdaP = lambda
		lambda = L + (1-C)*f*sinAlpha*(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))
		iterLimit = iterLimit - 1
	}

	if iterLimit == 0 {
		return -1.0
	}

	uSq := cosSqAlpha * (a*a - b*b) / (b * b)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))
	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))
	s := b * A * (sigma - deltaSigma)
	s = s * math.Pow10(-3)
	s = (Round(s, 0.5, 3))
	return s
}
