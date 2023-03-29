package ports

import (
	"strconv"
	"strings"
  "errors"
)

// This is a closure that accepts a port ranges in string representation like this:
// `3000-8000,60000-65000` and returns a function upon calling which it returns two
// values: next port number and an error in case if no more ports left from the
// ranges of ports supplied.
func InitPortIterator(rangeStr string) func() (int, error) {
	rangeSlice := strings.Split(rangeStr, ",")
	i := 0
	ranges := make([][]int, len(rangeSlice))
	for i, r := range rangeSlice {
		rangeVals := strings.Split(strings.TrimSpace(r), "-")
		start, _ := strconv.Atoi(rangeVals[0])
		end, _ := strconv.Atoi(rangeVals[1])

		if start > end {
			panic("😱 Invalid range supplied!")
		}

		ranges[i] = []int{start, end}

	}
	i = ranges[0][0]
	j := 0

  // Closure captures range variables
	return func() (int, error) {
		if i > ranges[j][1] {
			j++
			if j >= len(ranges) {
        j--
				return 0, errors.New("no more ports left")
			}
			i = ranges[j][0]
		}
		val := i
		i++
		return val, nil
	}
}
