package cursor

import (
	"encoding/base64"
	"encoding/json"
)

func Encode(data any) (string, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dataJson), nil
}

func Decode(in string, to any) error {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(decoded, to)
}
