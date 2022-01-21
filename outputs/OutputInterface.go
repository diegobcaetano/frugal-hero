package outputs

import "io"

type OutputInterface interface {
	io.Reader
	Write()
}
