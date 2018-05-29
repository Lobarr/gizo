package helpers

import "encoding/json"

func Deserialize(b []byte, d interface{}) error {
	err := json.Unmarshal(b, d)
	return err
}

func Serialize(d interface{}) ([]byte, error) {
	bytes, err := json.Marshal(d)
	return bytes, err
}
