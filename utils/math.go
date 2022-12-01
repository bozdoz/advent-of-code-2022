package utils

type numeric interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Max[T numeric](nums ...T) T {
	max := nums[0]

	for _, val := range nums {
		if val > max {
			max = val
		}
	}

	return max
}
