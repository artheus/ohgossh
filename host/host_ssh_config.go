package host

import (
	"fmt"
	"github.com/artheus/ohgossh/prompt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"os"
)

func (h *Host) SSHClientConfig() (sshConf *ssh.ClientConfig) {
	var hostKeyCallback ssh.HostKeyCallback

	if h.InsecureIgnoreHostKey {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	} else {
		hostKeyCallback = h.certCheckerCallback
	}

	sshConf = &ssh.ClientConfig{}
	sshConf.SetDefaults()

	sshConf.User = h.User
	sshConf.HostKeyCallback = hostKeyCallback
	sshConf.BannerCallback = h.bannerCallback()
	sshConf.Timeout = h.Timeout

	return sshConf
}

func (h *Host) hostKeyCallback() func(string, net.Addr, ssh.PublicKey) error {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) (err error) {
		var answer bool
		var marshaledHostKey = string(ssh.MarshalAuthorizedKey(key))

		if answer, err = prompt.YesNo(fmt.Sprintf("Allow host %s (%s) public key\n%s", hostname, remote.String(), marshaledHostKey)); err != nil {
			return errors.Wrap(err, "failed to prompt for host key approval")
		} else if answer == true {
			var knownHostsLine = knownhosts.Line([]string{hostname}, key)
			var knownHostsFile *os.File

			if knownHostsFile, err = os.OpenFile(h.Config.KnownHostsFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640); err != nil {
				return err
			}

			if _, err = knownHostsFile.WriteString(knownHostsLine); err != nil {
				return err
			}

			if err = knownHostsFile.Close(); err != nil {
				logrus.Warnf("failed to close known_hosts file: %s", err)
			}

			// Note: reload the known_hosts file
			if h.certCheckerCallback, err = knownhosts.New(h.Config.KnownHostsFilename); err != nil {
				return err
			}

			return nil
		}

		return errors.New("user did not approve host key")
	}
}

func (h *Host) bannerCallback() func(string) error {
	return func(message string) error {
		fmt.Println(message)

		return nil
	}
}
