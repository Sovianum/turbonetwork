package configuration

func GetVariants(limits []int) [][]int {
	total := 1
	for _, limit := range limits {
		total *= limit
	}

	result := make([][]int, total)
	dimensions := getDimensions(limits)

	for i := 0; i != total; i++ {
		result[i] = make([]int, len(limits))

		for j := range limits{
			incNum := 0
			if j == 0 {
				incNum = i
			} else {
				incNum = i / dimensions[j - 1]
			}
			result[i][j] = incNum % limits[j]
		}
	}
	return result
}

func getDimensions(limits []int) []int {
	total := 1
	result := make([]int, len(limits))

	for i, limit := range limits {
		total *= limit
		result[i] = total
	}
	return result
}
