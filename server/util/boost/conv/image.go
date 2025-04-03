package conv

import (
	"math"
	"onij/boost/exp"
)

// mm to px, default 72 dpi
func MillimeterToPixel(mm float64) int32 {
	if mm <= 0 {
		return exp.Zero[int32]()
	}

	inch := mm / 25.4
	return int32(math.Round(inch * 72))
}

func MillimeterToDPIPixel(mm, dpi float64) float64 {
	if mm <= 0 {
		return exp.Zero[float64]()
	}

	inch := mm / 25.4
	px := inch * dpi
	return px
}
