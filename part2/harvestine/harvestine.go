package harvestine

import "math"

// Square calculates the square of a number
func Square(A float64) float64 {
	Result := A * A
	return Result
}

// RadiansFromDegrees converts degrees to radians
func RadiansFromDegrees(Degrees float64) float64 {
	Result := 0.01745329251994329577 * Degrees
	return Result
}

// ReferenceHaversine calculates the Haversine distance between two points on Earth
// NOTE: EarthRadius is generally expected to be 6372.8
func ReferenceHaversine(X0, Y0, X1, Y1, EarthRadius float64) float64 {
	/* NOTE: This is not meant to be a "good" way to calculate the Haversine distance.
	   Instead, it attempts to follow, as closely as possible, the formula used in the real-world
	   question on which these homework exercises are loosely based.
	*/

	lat1 := Y0
	lat2 := Y1
	lon1 := X0
	lon2 := X1

	dLat := RadiansFromDegrees(lat2 - lat1)
	dLon := RadiansFromDegrees(lon2 - lon1)
	lat1 = RadiansFromDegrees(lat1)
	lat2 = RadiansFromDegrees(lat2)

	a := Square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*Square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a))

	Result := EarthRadius * c

	return Result
}
