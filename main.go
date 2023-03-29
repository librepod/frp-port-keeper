package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
  "github.com/librepod/frp-manager/ports"
  "github.com/librepod/frp-manager/store"
)

var db = store.CreateStore()
var nextPort = ports.InitPortIterator("1-3")

type Body struct {
	MachineID string `json:"machineid" binding:"required"`
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

	r.POST("/register", func(c *gin.Context) {

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

    r := new(store.Record)

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

      val := store.Record{
        Port: p,
        IP:   c.ClientIP(),
      }

      err = db.Set(body.MachineID, val)
      if err != nil {
        fmt.Printf("Error setting value: %v.\n", err)
        panic(err)
      }
      remotePort = p
		}

    c.JSON(http.StatusAccepted, gin.H{"status": "ok", "remotePort": remotePort})
	})

	// Get user value
	r.GET("/registrations/:machineId", func(c *gin.Context) {
		machineId := c.Params.ByName("machineId")
		retrievedVal := new(store.Record)
		data, err := db.Get(machineId, retrievedVal)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"machineId": machineId, "value": data})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
