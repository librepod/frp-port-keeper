package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"main/ports"

	"github.com/gin-gonic/gin"
)

// NewProxy sample payload body
// {
// 	"version":"0.1.0",
// 	"op":"NewProxy",
// 	"content":{
// 		"user":{
// 			"user":"alexmachine",
// 			"metas":null,
// 			"run_id":"xxxxxxxxxxxxxxxx"
// 		},
// 	"proxy_name":"xxxxxxxxxxx.librepod-relay",
// 	"proxy_type":"udp"
// 	}
// }

type Request struct {
	Version string  `json:"version"`
	Op      string  `json:"op"`
	Content Content `json:"content"`
}

type Content struct {
	User      User   `json:"user"`
	ProxyName string `json:"proxy_name"`
	ProxyType string `json:"proxy_type"`
	ProxyPort int    `json:"remot_port,omitempty"`
}

type User struct {
	User  string                 `json:"user"`
	Metas map[string]interface{} `json:"metas,omitempty"`
	RunID string                 `json:"run_id"`
}

type Response struct {
	Reject       bool    `json:"reject,omitempty"`
	RejectReason string  `json:"reject_reason,omitempty"`
	Unchange     bool    `json:"unchange"`
	Content      Content `json:"content,omitempty"`
}

func PortRegistrationsHandler(c *gin.Context) {
	fmt.Printf("Query params: %+v\n", c.Request.URL.Query())
	fmt.Printf("Body payload: %+v\n", c.Copy().Request.Body)
	// Retrieve the request body
	var requestBody Request
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		fmt.Println("error trying to bind request body")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "VALIDATEERR",
				"message": "Invalid inputs. Please check your payload",
			},
		)
		return
	}

	noMorePortsBody := gin.H{"reject": true, "reject_reason": "no more ports left"}
	defaultResBody := gin.H{"reject": false, "unchange": true}

	op := c.Query("op")
	fmt.Printf("🔅 Got request Operation: %s\n", op)

	// Pretty print the request body to console
	requestBodyBytes, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
	}
	fmt.Println("Request body: ", string(requestBodyBytes))

	if op != "NewProxy" {
		fmt.Printf("returning default response...\n")
		c.JSON(http.StatusOK, defaultResBody)
		return
	}

	fmt.Println("👋 Aloha!!! We got NewProxy operation!")

	p, err := ports.GetFreePort()
	if err != nil && err.Error() == "no more free ports left" {
		fmt.Println("No more ports left")
		c.JSON(http.StatusForbidden, noMorePortsBody)
		return
	}

	var res = Response{
		Unchange: false,
	// Save the alocated port to DB
	// var db = store.CreateStore()
	// db.Set()

		Content: Content{
			User: requestBody.Content.User,
			ProxyName: requestBody.Content.ProxyName,
			ProxyType: requestBody.Content.ProxyType,
			ProxyPort: p,
		},
	}

	responseBodyBytes, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println("Allowing NewProxy connection...")
	fmt.Println("Response body: ", string(responseBodyBytes))

	c.JSON(http.StatusOK, res)
}
