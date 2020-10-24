package common

import (
	"encoding/json"
	"io"
)

// JSONConnIO represents abstraction of JSON connection communication
type JSONConnIO struct {
	encoder *json.Encoder
	decoder *json.Decoder
}

// NewJSONConnIO creates new instance of JSONConnIO
func NewJSONConnIO(readWriter io.ReadWriter) JSONConnIO {
	return JSONConnIO{encoder: json.NewEncoder(readWriter), decoder: json.NewDecoder(readWriter)}
}

// SendMessage send JSON message to client connection
func (conn JSONConnIO) SendMessage(message interface{}) error {
	err := conn.encoder.Encode(&message)

	return err
}

// GetMessage get message from client connection
func (conn JSONConnIO) GetMessage(message interface{}) error {
	err := conn.decoder.Decode(&message)

	return err
}
