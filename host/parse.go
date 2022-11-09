package host

import (
	"fmt"
	"github.com/artheus/ohgossh/config"
	regexp "github.com/gijsbers/go-pcre"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func Parse(hostURL *url.URL, conf *config.Config) (host *Host, err error) {
	var hostnameMatcher *regexp.Matcher

	if host, err = NewHost(hostURL, conf); err != nil {
		return nil, err
	}

	logrus.Debugf("hostURL: %s", hostURL)
	logrus.Debugf("conf: %+v", conf)

	copyHostParams(&conf.Defaults, &host.HostParams)

	var hostname = hostURL.Hostname()

	host.Name = hostname

	for _, hostParams := range conf.Hosts {
		if hostParams.Name == "" && hostParams.Pattern == "" {
			hostParamsYaml, _ := yaml.Marshal(hostParams)
			return nil, errors.Errorf("host name or pattern are not provided for host %s", hostParamsYaml)
		}

		if hostParams.Name != "" && hostParams.Pattern != "" {
			hostParamsYaml, _ := yaml.Marshal(hostParams)
			return nil, errors.Errorf("mutually exclusive parameters name and pattern are both provided for host %s", hostParamsYaml)
		}

		if hostParams.Name != "" && hostname == hostParams.Name {
			logrus.Debugf("hostname %s matched settings for %s", hostname, hostParams.Name)
			copyHostParams(&hostParams, &host.HostParams)
			continue
		}

		if len(hostParams.Aliases) != 0 {
			for _, alias := range hostParams.Aliases {
				if hostname == alias {
					logrus.Debugf("hostname %s matched settings for host alias %s", hostname, hostParams.Name)
					copyHostParams(&hostParams, &host.HostParams)
					continue
				}
			}
		}

		if hostParams.Pattern != "" {
			var rePattern regexp.Regexp

			if rePattern, err = regexp.Compile(
				fmt.Sprintf(
					"^%s$",
					hostParams.Pattern,
				),
				regexp.DOTALL&regexp.JAVASCRIPT_COMPAT&regexp.UTF8,
			); err != nil {
				return nil, errors.Wrapf(err, "unable to compile regexp: %s", hostParams.Pattern)
			}

			matcher := rePattern.MatcherString(hostname, regexp.NOTEMPTY)

			if matcher.Matches() {
				hostnameMatcher = matcher

				copyHostParams(&hostParams, &host.HostParams)

				if err = renderParams(&host.HostParams, hostname, hostnameMatcher); err != nil {
					return nil, errors.Wrap(err, "failed to render parameters")
				}
			}
		}
	}

	if host.Name != "" && host.Pattern != "" && host.Replace != "" {
		host.Name = ""
	}

	if err = renderParams(&host.HostParams, hostname, hostnameMatcher); err != nil {
		return nil, errors.Wrap(err, "failed to render parameters")
	}

	// override port number, if provided by command argument
	if hostURL.Port() != "" {
		var p uint64
		if p, err = strconv.ParseUint(hostURL.Port(), 10, 16); err != nil {
			return nil, errors.Wrap(err, "failed to parse port number")
		}
		host.Port = uint16(p)
	}

	// override username, if provided by command argument
	if hostURL.User.Username() != "" {
		host.User = hostURL.User.Username()
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.Debugf("Using following host config for host %s:", host.Name)
		ye := yaml.NewEncoder(os.Stdout)
		ye.SetIndent(2)
		defer ye.Close()
		ye.Encode(host.HostParams)
	}

	return host, err
}

func renderParams(host *config.HostParams, hostname string, hostnameMatcher *regexp.Matcher) (err error) {
	if host.Replace != "" && hostnameMatcher != nil {
		if host.Name, err = renderTemplate(hostname, host.Replace, hostnameMatcher); err != nil {
			return errors.Wrap(err, "failed to render hostname template")
		}

		host.Pattern = ""
		host.Replace = ""
	}

	if host.User, err = renderTemplate(hostname, host.User, hostnameMatcher); err != nil {
		return errors.Wrap(err, "failed to render username template")
	}

	if host.JumpHost, err = renderTemplate(hostname, host.JumpHost, hostnameMatcher); err != nil {
		return errors.Wrap(err, "failed to render jumphost template")
	}

	if host.IdentityFile, err = renderTemplate(hostname, host.IdentityFile, hostnameMatcher); err != nil {
		return errors.Wrap(err, "failed to render identity file template")
	}

	if strings.Contains(host.IdentityFile, "~") {
		if currentUser, err := user.Current(); err != nil {
			return errors.Wrap(err, "failed to get shell user")
		} else {
			host.IdentityFile = strings.Replace(host.IdentityFile, "~", currentUser.HomeDir, -1)
		}
	}

	if host.HttpProxy != nil && host.HttpProxy.Auth != nil {
		if host.HttpProxy.Auth.User, err = renderTemplate(hostname, host.HttpProxy.Auth.User, hostnameMatcher); err != nil {
			return errors.Wrap(err, "failed to render htto proxy user template")
		}

		if host.HttpProxy.Auth.Password, err = renderTemplate(hostname, host.HttpProxy.Auth.Password, hostnameMatcher); err != nil {
			return errors.Wrap(err, "failed to render http proxy template template")
		}
	}

	return
}

// TODO: Do this much better through reflect
func copyHostParams(from *config.HostParams, to *config.HostParams) {
	if from.Name != "" {
		to.Name = from.Name
	}

	if from.Port != 0 {
		to.Port = from.Port
	}

	if from.Pattern != "" {
		to.Pattern = from.Pattern
	}

	if from.JumpHost != "" {
		to.JumpHost = from.JumpHost
	}

	if from.HttpProxy != nil {
		to.HttpProxy = from.HttpProxy
	}

	if from.Replace != "" {
		to.Replace = from.Replace
	}

	if from.GssAPI != nil {
		to.GssAPI = from.GssAPI
	}

	if from.IdentityFile != "" {
		to.IdentityFile = from.IdentityFile
	}

	if len(from.PreferredAuth) != 0 {
		to.PreferredAuth = from.PreferredAuth
	}

	if from.User != "" {
		to.User = from.User
	}
}
