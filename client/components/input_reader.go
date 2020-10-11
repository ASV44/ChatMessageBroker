package components

import (
	"bufio"
	"os"
	"strings"
)

type InputReader struct {
	reader *bufio.Reader
}

func NewInputReader() InputReader {
	return InputReader{reader: bufio.NewReader(os.Stdin)}
}

func (inputReader InputReader) GetUserInput() string {
	data, _ := inputReader.reader.ReadString('\n')
	return strings.TrimSuffix(data, "\n")
}
