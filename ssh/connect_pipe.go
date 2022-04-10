package ssh

import (
	. "github.com/artheus/ohgossh/utils"
	"io"
	"os"
)

func pipeToShell(stdin io.Writer, stdout, stderr io.Reader) {
	go IgnoreErrIOCopy(stdin, os.Stdin)

	go IgnoreErrIOCopy(os.Stdout, stdout)

	go IgnoreErrIOCopy(os.Stderr, stderr)
}
