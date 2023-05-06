package main

import (
	"errors"
	"fmt"

	"gopkg.in/ini.v1"

	"main/ports"
	"main/server"
)

var nextPort func() (int, error)

// ServerConf contains information for a server service. It is
// recommended to use GetDefaultServerConf instead of creating this object
// directly, so that all unspecified fields have reasonable default values.
type ServerConf struct {
	// Original string.
	AllowPorts string `ini:"-" json:"-"`
}

func UnmarshalServerConfFromIni(source interface{}) (ServerConf, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:         false,
		InsensitiveSections: false,
		InsensitiveKeys:     false,
		IgnoreInlineComment: true,
		AllowBooleanKeys:    true,
	}, source)

	if err != nil {
		return ServerConf{}, err
	}

	s, err := f.GetSection("common")
	if err != nil {
		return ServerConf{}, err
	}

	common := ServerConf{}
	err = s.MapTo(&common)
	if err != nil {
		return ServerConf{}, err
	}

	// allow_ports
	allowPortStr := s.Key("allow_ports").String()
	if allowPortStr != "" {
		common.AllowPorts = allowPortStr
	} else {
		return ServerConf{}, errors.New("common.allow_ports not specified in config")
	}

	return common, nil
}

func init() {
	fmt.Println("🐔 Aloha from init func!")

	var commonSection, err = UnmarshalServerConfFromIni("./frps.ini")
	if err != nil {
		fmt.Println("We got error: ", err)
	}
	fmt.Println("We got allowPorts: ", commonSection.AllowPorts)

	ports.InitPortsIterator(commonSection.AllowPorts)

}

func main() {
	server.Start()
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
