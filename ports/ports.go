package ports

import (
	"errors"
	"fmt"
	"main/store"
	"strconv"
	"strings"
)

var nextPort func() (int, error)

func init() {
	// TODO: Pass port range from frps.ini config
	allowedPorts := "6007-6010,6013,6017"
	nextPort = createAllowPortsGenerator(allowedPorts)
}

func GetFreePort() (int, error) {
	fmt.Println("Getting a pree port")

	freePort := 0

	// Iterate through all the allowedPorts skeeping those that had been already
	// alocated to somebody (have records in DB)
	for p, err := nextPort(); err == nil; p, err = nextPort() {
		fmt.Printf("Trying port %s...\n", p)
		portRecord := store.PortRecord{}
		found, dbErr := store.DB.Get(strconv.Itoa(p), &portRecord)
		if dbErr != nil {
			fmt.Println("error occurred accessing db")
			panic(dbErr)
		}
		if !found {
			fmt.Println("Found a free port to use: ", p)
			freePort = p
			break
		}
	}

	if freePort == 0 {
		return 0, errors.New("no more free ports left")
	}

	// fmt.Println("🥳 We got a slice of available ports: ", availablePorts)
	return freePort, nil
}

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

// func initAvailablePortsIterator(availablePorts []int) func() (int, error) {
// 	i := 0
// 	return func() (int, error) {
// 		if i >= len(availablePorts) {
// 			return 0, errors.New("no more free ports left")
// 		}
// 		r := availablePorts[i]
// 		i++
// 		return r, nil
// 	}
// }


// func GetFreePort(userName string) (int, error) {
// 	r := new(store.MachineRecord)
// 	found, _ := db.Get(userName, r)
// 	fmt.Printf("Found record: %v.\n", found)
// 	var p = r.Port
//
// 	if !found {
// 		fmt.Println("Record does not exist. Setting it...")
// 		p, err := NextPort()
// 		fmt.Println("Got port: ", p)
// 		if err != nil {
// 			fmt.Println(err)
// 			return 0, err
// 		}
//
// 		date := time.Now().UTC()
// 		rec1 := store.MachineRecord{
// 			Port:      p,
// 			IP:        c.ClientIP(),
// 			CreatedAt: date,
// 		}
// 		rec2 := store.PortRecord{
// 			MachineID: body.MachineID,
// 			IP:        c.ClientIP(),
// 			CreatedAt: date,
// 		}
//
// 		err = db.Set(body.MachineID, rec1)
// 		if err != nil {
// 			fmt.Printf("Error setting value: %v.\n", err)
// 			panic(err)
// 		}
// 		err = db.Set(strconv.Itoa(p), rec2)
// 		if err != nil {
// 			fmt.Printf("Error setting value: %v.\n", err)
// 			panic(err)
// 		}
// 		remotePort = p
// 	}
//
// 	return 0, err
//
// 	return p, nil
// }
