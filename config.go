/*
JSON config parser
*/

package xrgrpc

import (
	"encoding/json"
	"os"
)

// DecodeJSONConfig marshalls the file into the object v
// XR self-signed certificates are issued for CN=ems.cisco.com
func DecodeJSONConfig(v interface{}, f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(v)
}
