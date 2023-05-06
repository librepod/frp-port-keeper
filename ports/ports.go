package ports

import (
	"errors"
	"fmt"
	"main/store"
	"strconv"
	"strings"
)

var NextPortGenerator func() (int, error)

// This is a closure that accepts port ranges in string representation like this:
// `3000-8000,60000-65000` and returns an iterator function which returns
// next port number and an error in case if no more ports left from the
// ranges of ports supplied.
func createAllowPortsGenerator(portsRange string) func() (int, error) {
	rangeSlice := strings.Split(portsRange, ",")
	i := 0
	ranges := make([][]int, len(rangeSlice))
	for i, r := range rangeSlice {
		if strings.Contains(r, "-") {
			rangeVals := strings.Split(strings.TrimSpace(r), "-")
			start, _ := strconv.Atoi(rangeVals[0])
			end, _ := strconv.Atoi(rangeVals[1])

			if start > end {
				panic("😱 invalid range supplied")
			}

			ranges[i] = []int{start, end}
		} else {
			port, _ := strconv.Atoi(r)
			ranges[i] = []int{port, port}
		}
	}

	i = ranges[0][0]
	j := 0

	// Closure captures range variables
	return func() (int, error) {
		if i > ranges[j][1] {
			j++
			if j >= len(ranges) {
				j--
				return 0, errors.New("no more free ports left")
			}
			i = ranges[j][0]
		}
		val := i
		i++
		return val, nil
	}
}

func initAvailablePortsIterator(availablePorts []int) func() (int, error) {
	i := 0
	return func() (int, error) {
		if i >= len(availablePorts) {
			return 0, errors.New("no more free ports left")
		}
		r := availablePorts[i]
		i++
		return r, nil
	}
}

func InitPortsIterator(portsRange string) func() (int, error) {
	fmt.Println("portsRange: ", portsRange)

	var nextPort = createAllowPortsGenerator(portsRange)
	var availablePorts = []int{}

	var db = store.CreateStore()
	defer db.Close()

	// Iterate through all the allowedPorts sorting out those that had been already
	// alocated to somebody (have records in DB)
	for p, err := nextPort(); err == nil; p, err = nextPort() {
		portRecord := store.PortRecord{}
		found, dbErr := db.Get(strconv.Itoa(p), &portRecord)
		if dbErr != nil {
			fmt.Println("error occurred accessing db")
			panic(dbErr)
		}
		if found {
			continue
		}

		availablePorts = append(availablePorts, p)
	}
	// fmt.Println("🥳 We got a slice of available ports: ", availablePorts)
	NextPortGenerator = initAvailablePortsIterator(availablePorts)

	return NextPortGenerator
}
