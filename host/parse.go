package host

import (
	"fmt"
	"github.com/artheus/ohgossh/config"
	"github.com/artheus/ohgossh/utils"
	regexp "github.com/gijsbers/go-pcre"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"net/url"
	"os/user"
	"strconv"
	"strings"
)

func Parse(hostURL *url.URL, conf *config.Config) (host *Host, err error) {
	defer utils.HandleError(&err)

	var hostnameMatcher *regexp.Matcher

	host, err = NewHost(hostURL, conf)
	utils.PanicOnError(err)

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

			rePattern, err = regexp.Compile(
				fmt.Sprintf(
					"^%s$",
					hostParams.Pattern,
				),
				regexp.DOTALL&regexp.JAVASCRIPT_COMPAT&regexp.UTF8,
			)
			utils.PanicOnError(errors.Wrapf(err, "unable to compile regexp: %s", hostParams.Pattern))

			matcher := rePattern.MatcherString(hostname, regexp.NOTEMPTY)

			if matcher.Matches() {
				hostnameMatcher = matcher

				copyHostParams(&hostParams, &host.HostParams)

				utils.PanicOnError(
					renderParams(&host.HostParams, hostname, hostnameMatcher),
				)
			}
		}
	}

	if host.Name != "" && host.Pattern != "" && host.Replace != "" {
		host.Name = ""
	}

	utils.PanicOnError(
		renderParams(&host.HostParams, hostname, hostnameMatcher),
	)

	// override port number, if provided by command argument
	if hostURL.Port() != "" {
		var p uint64
		p, err = strconv.ParseUint(hostURL.Port(), 10, 16)
		utils.PanicOnError(err)
		host.Port = uint16(p)
	}

	// override username, if provided by command argument
	if hostURL.User.Username() != "" {
		host.User = hostURL.User.Username()
	}

	return host, err
}

func renderParams(host *config.HostParams, hostname string, hostnameMatcher *regexp.Matcher) (err error) {
	defer utils.HandleError(&err)

	if host.Replace != "" && hostnameMatcher != nil {
		host.Name, err = renderTemplate(hostname, host.Replace, hostnameMatcher)
		utils.PanicOnError(err)

		host.Pattern = ""
		host.Replace = ""
	}

	host.User, err = renderTemplate(hostname, host.User, hostnameMatcher)
	utils.PanicOnError(err)

	host.JumpHost, err = renderTemplate(hostname, host.JumpHost, hostnameMatcher)
	utils.PanicOnError(err)

	host.IdentityFile, err = renderTemplate(hostname, host.IdentityFile, hostnameMatcher)
	utils.PanicOnError(err)

	if strings.Contains(host.IdentityFile, "~") {
		currentUser, err := user.Current()
		utils.PanicOnError(err)

		host.IdentityFile = strings.Replace(host.IdentityFile, "~", currentUser.HomeDir, -1)
	}

	if host.HttpProxy != nil && host.HttpProxy.Auth != nil {
		host.HttpProxy.Auth.User, err = renderTemplate(hostname, host.HttpProxy.Auth.User, hostnameMatcher)
		utils.PanicOnError(err)

		host.HttpProxy.Auth.Password, err = renderTemplate(hostname, host.HttpProxy.Auth.Password, hostnameMatcher)
		utils.PanicOnError(err)
	}

	return
}

// TODO: Do this much better through reflect
func copyHostParams(from *config.HostParams, to *config.HostParams) {
	if from.Name != "" {
		to.Name = from.Name
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
