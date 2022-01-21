package outputs

import (
	"os"
)

type StdoutOutput struct {
	buffer []byte
}

func (stdo *StdoutOutput) Read(p []byte) (n int, err error) {
	stdo.buffer = append(stdo.buffer, p...)
	return 1, nil
}

func (stdo *StdoutOutput) Write() {
	_, err := os.Stdout.Write(stdo.buffer)

	if err != nil {
		print(err)
	}
}