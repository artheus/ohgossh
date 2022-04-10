package host

import (
	"fmt"
	"github.com/artheus/ohgossh/config"
	"github.com/artheus/ohgossh/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Host struct {
	config.HostParams

	URL    *url.URL
	Config *config.Config

	certCheckerCallback ssh.HostKeyCallback
}

func NewHost(hostURL *url.URL, conf *config.Config) (*Host, error) {
	var hkc, err = knownhosts.New(conf.KnownHostsFilename)

	if err != nil {
		return nil, err
	}

	return &Host{
		HostParams:          config.DefaultHostParams(),
		URL:                 hostURL,
		Config:              conf,
		certCheckerCallback: hkc,
	}, err
}

func (h *Host) GetJumpHost() (_ *Host, err error) {
	utils.HandleError(&err)

	if h.JumpHost == "" {
		return nil, nil
	}

	var jumpHostURL *url.URL

	if !strings.HasPrefix(h.JumpHost, "ssh://") {
		h.JumpHost = fmt.Sprintf("ssh://%s", h.JumpHost)
	}

	jumpHostURL, err = url.Parse(h.JumpHost)

	return Parse(jumpHostURL, h.Config)
}

func (h *Host) Addr() string {
	var port = "22"

	if h.Port != 0 {
		port = strconv.FormatUint(uint64(h.Port), 10)
	}

	return net.JoinHostPort(h.Name, port)
}
