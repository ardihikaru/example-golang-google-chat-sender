package json

import (
	"encoding/json"
	"fmt"
)

func BuildJsonByteStr(docs interface{}) (*[]byte, error) {
	var err error

	// builds json bytes string
	var jsonBytesStr []byte
	jsonBytesStr, err = json.Marshal(docs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert JSON object to string: %s", err.Error())
	}

	return &jsonBytesStr, nil
}
