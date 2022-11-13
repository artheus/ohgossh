package ssh

import (
	"fmt"
	"github.com/artheus/ohgossh/host"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"syscall"
)

func Connect(host *host.Host) (client *ssh.Client, err error) {
	jumpHost, _ := host.GetJumpHost()

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		blue := color.New(color.FgBlue).SprintFunc()
		var logString = fmt.Sprintf("connecting to host %s", blue(host.Addr()))

		if jumpHost != nil {
			logString = fmt.Sprintf("%s, through jumphost %s", logString, blue(jumpHost.Name))
		}

		logrus.Debug(logString)
	}

	client, err = Dial(host, jumpHost)
	if err != nil {
		if jumpHost != nil {
			return nil, errors.Wrapf(err, "failed to dial host %s, jumphost: %s", host.Addr(), jumpHost.Addr())
		} else {
			return nil, errors.Wrapf(err, "failed to dial host %s", host.Addr())
		}
	}

	return client, err
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
