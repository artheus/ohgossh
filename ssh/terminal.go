package ssh

import (
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"os"
	"syscall"
)

var modes = ssh.TerminalModes{
	ssh.ECHO:          1,
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
}

func windowSize() (tw int, th int) {
	tw, th, err := term.GetSize(0)
	if err != nil {
		tw = 80
		th = 10
	}

	return tw, th
}

func isTTYAvailable() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

func RequestPTYIfTTYAvailable(session *ssh.Session) (restoreFunc func(), err error) {
	// Request PTY and clear screen, if shell is a TTY
	if isTTYAvailable() {
		tw, th := windowSize()

		var oldStdinState *term.State

		// make terminal raw, to send everything on stdin to ssh session
		oldStdinState, err = term.MakeRaw(syscall.Stdin)
		if err != nil {
			return nil, errors.Wrap(err, "failed with term.MakeRaw")
		}

		restoreFunc = func(restoreState *term.State) func() {
			return func() {
				restoreTermState(oldStdinState)
			}
		}(oldStdinState)

		err = session.RequestPty("xterm", th, tw, modes)
		if err != nil {
			return nil, errors.Wrap(err, "unable to request pty for ssh session")
		}

		clearShell()
	}

	if restoreFunc == nil {
		restoreFunc = func() {}
	}

	return restoreFunc, nil
}
