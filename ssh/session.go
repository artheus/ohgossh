package ssh

import (
	. "github.com/artheus/ohgossh/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"strings"
)

func RunSession(client *ssh.Client, command []string) (err error) {
	var session *ssh.Session

	session, err = client.NewSession()
	if err != nil {
		return errors.Wrap(err, "failed to open ssh session")
	}

	inPipe, outPipe, errPipe, err := getSessionPipes(session)
	if err != nil {
		return err
	}

	if inPipe != nil && outPipe != nil && errPipe != nil {
		pipeToShell(inPipe, outPipe, errPipe)
	}

	var restoreTerm func()

	if restoreTerm, err = RequestPTYIfTTYAvailable(session); err != nil {
		return errors.Wrap(err, "PTY request failed")
	}

	defer restoreTerm()

	if len(command) == 0 {
		// NOTE: Request shell if no command is provided
		err = session.Shell()
		if err != nil {
			return errors.Wrap(err, "failed to start remote shell")
		}
	} else {
		// NOTE: Run command on the remote host if command is provided
		err = session.Start(strings.Join(command, " "))
		if err != nil {
			return errors.Wrap(err, "failed to run remote command")
		}
	}

	return errors.WithStack(session.Wait())
}

func getSessionPipes(session *ssh.Session) (stdin io.Writer, stdout, stderr io.Reader, err error) {
	stdin, err = session.StdinPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to get ssh session stdin pipe")
	}

	stdout, err = session.StdoutPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to get ssh session stdout pipe")
	}

	stderr, err = session.StderrPipe()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "unable to get ssh session stderr pipe")
	}

	return
}

func pipeToShell(stdin io.Writer, stdout, stderr io.Reader) {
	go IgnoreErrIOCopy(stdin, os.Stdin)

	go IgnoreErrIOCopy(os.Stdout, stdout)

	go IgnoreErrIOCopy(os.Stderr, stderr)
}
