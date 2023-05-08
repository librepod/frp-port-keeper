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

	// Validate the request payload
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

	op := c.Query("op")
	fmt.Printf("Got %s request operation\n", op)

	// Pretty print the request body to console
	requestBodyBytes, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
	}
	fmt.Println("Request body: ", string(requestBodyBytes))

	if op != "NewProxy" {
		fmt.Printf("returning default response...\n")
		c.JSON(http.StatusOK, gin.H{"reject": false, "unchange": true})
		return
	}

	fmt.Println("Processing NewProxy operation...")

	p, err := ports.GetFreePort(requestBody.Content.User.User)
	if err != nil {
		fmt.Printf("Error ocurred getting free port: %+v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"reject": true, "reject_reason": err.Error()})
		return
	}

	var res = Response{
		Unchange: false,
		Content: Content{
			User:      requestBody.Content.User,
			ProxyName: requestBody.Content.ProxyName,
			ProxyType: requestBody.Content.ProxyType,
			ProxyPort: p,
		},
	}

	responseBodyBytes, _ := json.MarshalIndent(res, "", "  ")
	fmt.Printf("Allowing connection with response payload: %+v\n", string(responseBodyBytes))
	c.JSON(http.StatusOK, res)
}
