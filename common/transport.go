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
	return conn.encoder.Encode(&message)
}

// GetMessage get message from client connection
func (conn JSONConnIO) GetMessage(message interface{}) error {
	return conn.decoder.Decode(&message)
}
