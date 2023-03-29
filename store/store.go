package store

import (
	"encoding/json"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/file"
)

type Record struct {
	Port int
	IP   string
}


func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func CreateStore() gokv.Store {
	options := file.DefaultOptions
	store, err := file.NewStore(options)
	if err != nil {
		panic(err)
	}
	return store
}
