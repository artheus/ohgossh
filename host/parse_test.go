package host

import (
	"github.com/artheus/ohgossh/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopyHostParams(t *testing.T) {
	fromHostParams := new(config.HostParams)
	toHostParams := new(config.HostParams)

	fromHostParams.Name = "this-is-my.name.org"
	fromHostParams.JumpHost = "jumpyard"
	fromHostParams.HttpProxy = &config.HttpProxyParams{
		Host: "sooopar-proxy1.local",
		Port: 8080,
	}

	copyHostParams(fromHostParams, toHostParams)

	assert.Equal(t, fromHostParams.Name, toHostParams.Name)
	assert.Equal(t, fromHostParams.JumpHost, toHostParams.JumpHost)
	assert.Equal(t, fromHostParams.HttpProxy, toHostParams.HttpProxy)
}
