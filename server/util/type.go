package util

func IntToBool[R int | int16 | int32 | int64](v R) bool {
	return v > 0
}
func BoolToInt16(v bool) int16 {
	return IfElse(v, int16(1), int16(0))
}
