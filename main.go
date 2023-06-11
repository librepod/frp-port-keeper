package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"

	"main/ports"
	"main/server"
	"main/store"
)

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
		fmt.Println("⚠ common.allow_ports not specified in config, falling back to 1000-65535 port range")
		common.AllowPorts = "1000-65535"
	}

	return common, nil
}

func init() {
	fmt.Println("🐔 Initializing the plugin...")

	// Get the frps.ini file path from the environment variable
	frpsIniPath := os.Getenv("FRPS_INI_PATH")

	// If the environment variable is empty or not set
	if frpsIniPath == "" {
		// Use the default directory
		frpsIniPath = "./frps.ini"
	}

	// Check if frps.ini exists
	if _, err := os.Stat(frpsIniPath); os.IsNotExist(err) {
		panicMsg, _ := fmt.Printf("frps.ini does not exist at path %s; move the frp-port-keeper binary to the same folder where the frps.ini located and call frp-port-keeper from there.", frpsIniPath)
		panic(panicMsg)
	}

	var commonSection, err = UnmarshalServerConfFromIni(frpsIniPath)
	if err != nil {
		fmt.Println("got error: ", err)
	}
	fmt.Println("got allow_ports value: ", commonSection.AllowPorts)

	ports.InitPortsGenerator(commonSection.AllowPorts)
}

func main() {
	defer store.DB.Close()
	server.Start()
}
