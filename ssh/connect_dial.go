package ssh

import (
	"github.com/artheus/ohgossh/host"
	"github.com/artheus/ohgossh/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
)

func Dial(host *host.Host, jumpHost *host.Host) (client *ssh.Client, err error) {
	if jumpHost != nil {
		var conn net.Conn
		var jumpHostClient *ssh.Client

		jumpHostClientConfig := jumpHost.SSHClientConfig()

		logrus.Debugf("ssh dial jumphost %s", jumpHost.Addr())
		if jumpHostClient, err = ssh.Dial("tcp", jumpHost.Addr(), jumpHostClientConfig); err != nil {
			return nil, err
		}

		logrus.Debugf("ssh client dial jumphost %s", jumpHost.Addr())
		if conn, err = jumpHostClient.Dial("tcp", host.Addr()); err != nil {
			return nil, err
		}

		return tryConnect(host, func(sshConf *ssh.ClientConfig) (client *ssh.Client, err error) {
			var clientConn ssh.Conn
			var newChan <-chan ssh.NewChannel
			var reqChan <-chan *ssh.Request

			if clientConn, newChan, reqChan, err = ssh.NewClientConn(conn, host.Addr(), sshConf); err != nil {
				return nil, errors.Wrap(err, "failed to create client connection")
			}

			client = ssh.NewClient(clientConn, newChan, reqChan)

			return
		})
	}

	return tryConnect(
		host,
		func(sshConf *ssh.ClientConfig) (client *ssh.Client, err error) {
			return ssh.Dial("tcp", host.Addr(), sshConf)
		},
	)
}

type connectFunc func(sshConf *ssh.ClientConfig) (*ssh.Client, error)

func tryConnect(host *host.Host, c connectFunc) (client *ssh.Client, err error) {
	utils.HandleError(&err)

	var authMethod ssh.AuthMethod
	var sshConf = host.SSHClientConfig()

	logrus.Debugf("Attempting auth methods %v", host.PreferredAuth)

	for _, methodStr := range host.PreferredAuth {
		authMethod, err = authMethodFor(methodStr, host)
		utils.PanicOnError(err)

		sshConf.Auth = []ssh.AuthMethod{authMethod}

		if client, err = c(sshConf); err != nil {
			continue
		}

		return
	}

	return nil, errors.Errorf("tried all auth methods in %v, of which none succeeded", host.PreferredAuth)
}
