package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/librepod/frp-manager/ports"
	"github.com/librepod/frp-manager/store"
)

var nextPort func() (int, error)

type Body struct {
	MachineID string `json:"machineid" binding:"required"`
}

func init() {
  fmt.Println("🐔 Aloha from init func!")
  var availablePorts = []int{}
  var next = ports.InitAllPortsIterator("6000-6005")
  var db = store.CreateStore()
  defer db.Close()

  fmt.Println("starting for loop...")

  for p, err := next(); err == nil; p, err = next() {
    fmt.Println("port: ", p)
    portRecord := store.PortRecord{}
    found, dbErr := db.Get(strconv.Itoa(p), &portRecord)
    if dbErr != nil {
      fmt.Println("Error occurred accessing db")
      panic(dbErr)
    }
    if found {
      continue
    }
    availablePorts = append(availablePorts, p)
  } 
  fmt.Println("🥳 We got a slice of available ports: ", availablePorts)
  nextPort = ports.InitAvailablePortsIterator(availablePorts)
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	gin.ForceConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/port-registrations", func(c *gin.Context) {

		body := Body{}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{
					"error":   "VALIDATEERR",
					"message": "Invalid inputs. Please check your inputs",
				},
			)
			return
		}
		fmt.Printf("body: %v.\n", body)

		r := new(store.MachineRecord)

		var db = store.CreateStore()
		defer db.Close()
		found, _ := db.Get(body.MachineID, r)

		fmt.Printf("Found record: %v.\n", found)

		var remotePort = r.Port

		if !found {
			fmt.Println("Record does not exist. Setting it...")
			p, err := nextPort()
			fmt.Println("Got port: ", p)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
				return
			}

      date := time.Now().UTC()
			rec1 := store.MachineRecord{
				Port:      p,
				IP:        c.ClientIP(),
				CreatedAt: date,
			}
			rec2 := store.PortRecord{
				MachineID: body.MachineID,
				IP:        c.ClientIP(),
				CreatedAt: date,
			}

			err = db.Set(body.MachineID, rec1)
			if err != nil {
				fmt.Printf("Error setting value: %v.\n", err)
				panic(err)
			}
			err = db.Set(strconv.Itoa(p), rec2)
			if err != nil {
				fmt.Printf("Error setting value: %v.\n", err)
				panic(err)
			}
			remotePort = p
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "ok", "remotePort": remotePort})
	})

	// // Get user value
	// r.GET("/port-registrations/:machineId", func(c *gin.Context) {
	// 	machineId := c.Params.ByName("machineId")
	// 	retrievedVal := new(store.Record)
	// 	var db = store.CreateStore()
	// 	defer db.Close()
	// 	data, err := db.Get(machineId, retrievedVal)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	c.JSON(http.StatusOK, gin.H{"machineId": machineId, "value": data})
	// })

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
