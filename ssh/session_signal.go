package ssh

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"os"
	"os/signal"
	"syscall"
)

var notifySignals = []os.Signal{
	syscall.SIGABRT,
	syscall.SIGALRM,
	syscall.SIGFPE,
	syscall.SIGHUP,
	syscall.SIGILL,
	syscall.SIGINT,
	syscall.SIGKILL,
	syscall.SIGPIPE,
	syscall.SIGQUIT,
	syscall.SIGSEGV,
	syscall.SIGTERM,
	syscall.SIGUSR1,
	syscall.SIGUSR2,
}

var signalMap = map[os.Signal]ssh.Signal{
	syscall.SIGABRT: ssh.SIGABRT,
	syscall.SIGALRM: ssh.SIGALRM,
	syscall.SIGFPE:  ssh.SIGFPE,
	syscall.SIGHUP:  ssh.SIGHUP,
	syscall.SIGILL:  ssh.SIGILL,
	syscall.SIGINT:  ssh.SIGINT,
	syscall.SIGKILL: ssh.SIGKILL,
	syscall.SIGPIPE: ssh.SIGPIPE,
	syscall.SIGQUIT: ssh.SIGQUIT,
	syscall.SIGSEGV: ssh.SIGSEGV,
	syscall.SIGTERM: ssh.SIGTERM,
	syscall.SIGUSR1: ssh.SIGUSR1,
	syscall.SIGUSR2: ssh.SIGUSR2,
}

func forwardSignalsToSession(session *ssh.Session) (cancelCtx context.CancelFunc) {
	// Create context for signal forwarding
	ctx, cancelCtx := context.WithCancel(context.Background())

	// Listen for os signals
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan)

	// Forward signals to SSH session
	go func(ctx context.Context, cancelCtx context.CancelFunc, signalChan <-chan os.Signal, session *ssh.Session) {
		select {
		case <-ctx.Done():
			return
		case sig := <-signalChan:

			logrus.Infof("got signal: %+v", sig)

			switch sig {
			case syscall.SIGQUIT:
				_ = session.Close()
			}

			if mappedSig, ok := signalMap[sig]; ok {
				logrus.Infof("sending signal: %+v", mappedSig)
				if err := session.Signal(mappedSig); err != nil {
					logrus.Warnf("failed to send signal %s to ssh session: %+v", sig.String(), err)
				}
			} else {
				logrus.Info("could not map signal, ignoring")
			}
		}
	}(ctx, cancelCtx, signalChan, session)

	return cancelCtx
}
