package utils

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func NewInt64(i int64) *int64 {
	return &i
}
