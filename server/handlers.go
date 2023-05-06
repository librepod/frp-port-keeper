package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"main/ports"
	"main/store"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Version string                 `json:"version"`
	Op      string                 `json:"op"`
	Content map[string]interface{} `json:"content"`
}

type Response struct {
	Reject       bool                   `json:"reject"`
	RejectReason string                 `json:"reject_reason"`
	Unchange     bool                   `json:"unchange"`
	Content      map[string]interface{} `json:"content"`
}

type Body struct {
	MachineID string `json:"machineid"`
	// MachineID string `json:"machineid" binding:"required"`
}

func PortRegistrationsHandler(c *gin.Context) {
	// fmt.Printf("Body payload: %+v\n", c.Copy().Request.Body)
	// Retrieve the request body
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	op := c.Query("op")
	fmt.Printf("🔅 Got request Operation: %s\n", op)

	// bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
	// jsonBody := string(bodyAsByteArray)

	// Pretty print the request body to console
	requestBodyBytes, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
	}
	// prettyJson, _ := json.MarshalIndent(jsonBody, "", "\t")
	// fmt.Printf("%s\n", jsonBody)

	switch op {
	case "NewProxy":
		fmt.Println("Aloha!!! We got NewProxy!")
		fmt.Println("Request body:", string(requestBodyBytes))
		c.JSON(http.StatusForbidden, gin.H{"reject": true, "reject_reason": "invalid user"})
	default:
		fmt.Printf("Default: %s\n", op)
		c.JSON(http.StatusOK, gin.H{"reject": false, "unchange": true})
	}

	// version := c.Query("version")

	// c.JSON(http.StatusForbidden, gin.H{"reject": true, "reject_reason": "invalid user"})
	// jsonData := []byte(`{"reject": false, "content": { "user": {"user": "kalich"}, "proxy_name": "docker-machine", "proxy_type": "udp", "remote_port": 6001}}`)
	// c.Data(http.StatusOK, "application/json", jsonData)
}

func HandlerOld(c *gin.Context) {

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
		p, err := ports.NextPort()
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
}

func NextPort() {
	panic("unimplemented")
}

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
