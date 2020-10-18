package services

import (
	"encoding/json"
	"fmt"
	"io"
)

type JsonConnIO struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewJsonConnIO(readWriter io.ReadWriter) JsonConnIO {
	return JsonConnIO{encoder: json.NewEncoder(readWriter), decoder: json.NewDecoder(readWriter)}
}

func (conn JsonConnIO) SendMessage(message interface{}) error {
	err := conn.encoder.Encode(&message)
	if err != nil {
		fmt.Println("Could not write message data ", err)
	}

	return err
}

func (conn JsonConnIO) GetMessage(message interface{}) error {
	err := conn.decoder.Decode(&message)
	if err != nil {
		fmt.Println("Could not decode message data ", err)
	}

	return err
}
