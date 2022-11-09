package config

import (
	"github.com/artheus/ohgossh/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"
)

type Config struct {
	Defaults HostParams   `yaml:"defaults,omitempty"`
	Hosts    []HostParams `yaml:"hosts,omitempty"`

	KnownHostsFilename string `yaml:"-"`
}

func LoadConfig(file string) (conf *Config, err error) {
	utils.HandleError(&err)

	confData, err := ioutil.ReadFile(file)
	utils.PanicOnError(err)

	conf = new(Config)

	utils.PanicOnError(
		yaml.Unmarshal(confData, conf),
	)

	return
}

type HostParams struct {
	Name                  string           `yaml:"name,omitempty"`
	Port                  uint16           `yaml:"port,omitempty"`
	User                  string           `yaml:"user,omitempty"`
	IdentityFile          string           `yaml:"identityFile,omitempty"`
	PreferredAuth         []string         `yaml:"preferredAuthentications,omitempty"`
	Aliases               []string         `yaml:"aliases,omitempty"`
	Pattern               string           `yaml:"pattern,omitempty"`
	Replace               string           `yaml:"replace,omitempty"`
	GssAPI                *GssAPIParams    `yaml:"gssApi,omitempty"`
	JumpHost              string           `yaml:"jumpHost,omitempty"`
	HttpProxy             *HttpProxyParams `yaml:"httpProxy,omitempty"`
	Timeout               time.Duration    `yaml:"timeout,omitempty"`
	InsecureIgnoreHostKey bool             `yaml:"ignoreHostKey,omitempty"`
}

type GssAPIParams struct {
	Enabled             bool `yaml:"enabled,omitempty"`
	DelegateCredentials bool `yaml:"delegateCredentials,omitempty"`
}

type HttpProxyParams struct {
	Host string               `yaml:"host"`
	Port uint16               `yaml:"port"`
	Auth *HttpProxyAuthParams `yaml:"auth"`
}

type HttpProxyAuthParams struct {
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
}
