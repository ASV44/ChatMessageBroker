package components

import (
	"bufio"
	"os"
	"strings"
)

// InputReader represents component which encapsulates logic of reading user input from standard input
type InputReader struct {
	reader *bufio.Reader
}

// NewInputReader creates new instance of InputReader
func NewInputReader() InputReader {
	return InputReader{reader: bufio.NewReader(os.Stdin)}
}

// GetUserInput read and returns user input from standard input
func (inputReader InputReader) GetUserInput() string {
	data, _ := inputReader.reader.ReadString('\n')
	return strings.TrimSuffix(data, "\n")
}
