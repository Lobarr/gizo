package helpers

import "encoding/json"

//Deserialize unmarshals from json
func Deserialize(b []byte, d interface{}) error {
	err := json.Unmarshal(b, d)
	return err
}

//Serialize marshals to json
func Serialize(d interface{}) ([]byte, error) {
	bytes, err := json.Marshal(d)
	return bytes, err
}
