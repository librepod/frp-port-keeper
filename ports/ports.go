package ports

import (
  "fmt"
	"strconv"
	"strings"
)

// func Next() func() int {
//   p := 0
//   fmt.Println("Generating next port...")
//
//   // Closure captures variable n
//   return func() int {
//     p += 1
//     return p
//   }
// }

func InitPortIterator(rangeStr string) func() (int, error) {
	rangeSlice := strings.Split(rangeStr, ",")
	i := 0
	ranges := make([][]int, len(rangeSlice))
	for i, r := range rangeSlice {
		rangeVals := strings.Split(strings.TrimSpace(r), "-")
		start, _ := strconv.Atoi(rangeVals[0])
		end, _ := strconv.Atoi(rangeVals[1])
		if start > end {
			panic("Invalid range")
		}
		ranges[i] = []int{start, end}

    fmt.Printf("ranges[i]: %v\n", ranges[i])
    fmt.Printf("ranges: %v\n", ranges)

	}
	i = ranges[0][0]
	j := 0

  // Closure captures range variables
	return func() (int, error) {
		if i > ranges[j][1] {
			j++
			if j >= len(ranges) {
				return 0, fmt.Errorf("😱 Whoa! No more ports left!")
			}
			i = ranges[j][0]
		}
		val := i
		i++
    fmt.Println("  Returning port: ", val)
		return val, nil
	}
}
