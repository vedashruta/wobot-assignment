package json

import "encoding/json"

func Encode(input any) (result []byte, err error) {
	result, err = json.Marshal(input)
	if err != nil {
		return
	}
	return
}

func Decode(input []byte, toAddress any) (err error) {
	err = json.Unmarshal(input, &toAddress)
	if err != nil {
		return
	}
	return
}
