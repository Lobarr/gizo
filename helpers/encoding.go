package helpers

import (
	"encoding/base64"
)

//Encode64 used to encode serilized block to base64 for writing to *.blk file
func Encode64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

//Decode64 used to decode base64
func Decode64(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}
