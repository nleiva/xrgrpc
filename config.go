package main

import (
	"encoding/json"
	"os"
)

// decodeJsonConfig marshalls the file into the object v
// XR self-signed certificates are issued for CN=ems.cisco.com
func decodeJSONConfig(v interface{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(v)
}
