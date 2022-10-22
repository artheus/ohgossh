package ssh

import (
	"fmt"
	"github.com/artheus/ohgossh/host"
	. "github.com/artheus/ohgossh/utils"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"strings"
	"syscall"
)

func Connect(host *host.Host, command []string) (err error) {
	defer HandleError(&err)

	jumpHost, _ := host.GetJumpHost()

	var client *ssh.Client

	blue := color.New(color.FgBlue).SprintFunc()
	var logString = fmt.Sprintf("connecting to host %s", blue(host.Name))

	if jumpHost != nil {
		logString = fmt.Sprintf("%s, through jumphost %s", logString, blue(jumpHost.Name))
	}

	// NOTE: won't need blue again..
	blue = nil

	logrus.Info(logString)

	client, err = Dial(host, jumpHost)
	PanicOnError(err)

	session, err := client.NewSession()
	PanicOnError(errors.Wrap(err, "failed to open ssh session"))

	defer Close(session)

	tw, th, err := term.GetSize(0)
	if err != nil {
		tw = 80
		th = 10
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// make terminal raw, to send everything on stdin to ssh session
	oldStdinState, err := term.MakeRaw(syscall.Stdin)
	PanicOnError(err)

	defer restoreTermState(oldStdinState)

	err = session.RequestPty("xterm", th, tw, modes)
	PanicOnError(errors.Wrap(err, "unable to request pty for ssh session"))

	errPipe, err := session.StderrPipe()
	PanicOnError(errors.Wrap(err, "unable to get ssh session stderr pipe"))

	outPipe, err := session.StdoutPipe()
	PanicOnError(errors.Wrap(err, "unable to get ssh session stdout pipe"))

	inPipe, err := session.StdinPipe()
	PanicOnError(errors.Wrap(err, "unable to get ssh session stdin pipe"))

	pipeToShell(inPipe, outPipe, errPipe)

	if len(command) == 0 {
		clearShell()

		err = session.Shell()
		PanicOnError(errors.Wrap(err, "failed to start remote shell"))
	} else {
		err = session.Start(strings.Join(command, " "))
		PanicOnError(errors.Wrap(err, "failed to run remote command"))
	}

	return session.Wait()
}

func restoreTermState(state *term.State) {
	if err := term.Restore(syscall.Stdin, state); err != nil {
		err = errors.Wrap(err, "failed to restore terminal state")

		logrus.Error(err)

		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			fmt.Printf("%+v\n", err)
		}
	}
}
