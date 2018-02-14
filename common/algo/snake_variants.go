package algo

const (
	up = iota
	down
	stay
)

func GetSnakeVariants(limits []int) [][]int {
	resultLength := 1
	for _, limit := range limits {
		resultLength *= limit
	}

	if resultLength == 0 {
		return [][]int{}
	}

	divisors := make([]int, len(limits))
	divisors[0] = 1
	curr := 1
	for i, limit := range limits[1:] {
		curr *= limit
		divisors[i + 1] = curr
	}

	snakes := make([][]int, len(limits))
	for i, divisor := range divisors {
		snakes[i] = getSnakeInts(0, limits[i] - 1, resultLength / divisor)
	}

	result := make([][]int, resultLength)
	for i := range result {
		item := make([]int, len(limits))
		for j := range item {
			item[j] = snakes[j][i / divisors[j]]
		}
		result[i] = item
	}
	return result
}

func getSnakeInts(start, end, len int) []int {
	result := make([]int, len)

	curr := start
	direction := up
	for i := 0; i != len; i++ {
		result[i] = curr

		if curr == start {
			if direction == down {
				direction = stay
			} else if direction == stay {
				direction = up
			}
		} else if curr == end {
			if direction == up {
				direction = stay
			} else if direction == stay {
				direction = down
			}
		}

		if direction == up && curr < end {
			curr++
		} else if direction == down && curr > start {
			curr--
		}
	}
	return result
}
