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

	// Check if frps.ini exists
	if _, err := os.Stat("./frps.ini"); os.IsNotExist(err) {
		panic("frps.ini does not exist; move the frp-port-keeper binary to the same folder where the frps.ini located and call frp-port-keeper from there.")
	}

	var commonSection, err = UnmarshalServerConfFromIni("./frps.ini")
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
