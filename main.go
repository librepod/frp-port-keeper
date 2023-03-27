package main

import (
	// "io/ioutil"
	"fmt"

	"encoding/json"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/philippgille/gokv"

	"github.com/philippgille/gokv/file"
  "github.com/librepod/frp-manager/ports"
)

type Record struct {
	Port int
	IP   string
}

// var db = make(map[string]string)
var db = createStore()

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func createStore() gokv.Store {
	options := file.DefaultOptions
	store, err := file.NewStore(options)
	if err != nil {
		panic(err)
	}
	return store
}

type Body struct {
	MachineID string `json:"machineid" binding:"required"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	gin.DisableConsoleColor()
	r := gin.Default()

  nextPort := ports.InitPortIterator("6000-7000")

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/register", func(c *gin.Context) {

		body := Body{}
		// using BindJson method to serialize body with struct
		// if err := c.BindJSON(&body); err != nil {
		//   c.AbortWithError(http.StatusBadRequest, err)
		//   return
		// }
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

    r := new(Record)

		found, _ := db.Get(body.MachineID, r)

    fmt.Printf("Found record: %v.\n", found)

    var remotePort = r.Port

		if !found {
      fmt.Println("Record does not exist. Setting it...")
      p, err := nextPort()
      fmt.Println("Got port: ", p)
      if err != nil {
        fmt.Println(err)
        c.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err})
      }

      val := Record{
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

		// c.JSON(http.StatusAccepted, &body)
    c.JSON(http.StatusOK, gin.H{"status": "ok", "remotePort": remotePort})

		// type MachineIDRequestBody struct {
		//   MachineID string
		// }

		// var requestBody MachineIDRequestBody

		// bodyAsByteArray, _ := ioutil.ReadAll(c.Request.Body)
		// jsonBody := string(bodyAsByteArray)
		// fmt.Println("JSON: \n" + jsonBody)

		// res, err := PrettyStruct(json)
		// if err != nil {
		//     log.Fatal(err)
		// }
		// fmt.Println("Here is your json: \n" + res)

		// if c.BindJSON(&requestBody) == nil {
		// 	val := record{
		// 		Port: 65001,
		// 		IP:   c.ClientIP(),
		// 	}
		// 	err := db.Set(requestBody.MachineID, val)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		//
		// 	c.JSON(http.StatusOK, gin.H{"status": "ok"})
		// }
	})

	// Get user value
	r.GET("/registrations/:machineId", func(c *gin.Context) {
		machineId := c.Params.ByName("machineId")
		retrievedVal := new(Record)
		data, err := db.Get(machineId, retrievedVal)
		if err != nil {
			panic(err)
		}
		// if !data {
		//   c.JSON(http.StatusOK, gin.H{"machineId": machineId, "value": "no value"})
		// }
		c.JSON(http.StatusOK, gin.H{"machineId": machineId, "value": data})
		// if err {
		// 	c.JSON(http.StatusOK, gin.H{"machineId": machineId, "value": data})
		// } else {
		// 	c.JSON(http.StatusOK, gin.H{"machineId": machineId, "status": "no value"})
		// }
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	// authorized.POST("admin", func(c *gin.Context) {
	// 	user := c.MustGet(gin.AuthUserKey).(string)
	//
	// 	// Parse JSON
	// 	var json struct {
	// 		Value string `json:"value" binding:"required"`
	// 	}
	//
	// 	if c.Bind(&json) == nil {
	// 		db[user] = json.Value
	// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// 	}
	// })

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
