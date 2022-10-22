package cmd

import (
	"fmt"
	"github.com/artheus/ohgossh/config"
	"github.com/artheus/ohgossh/host"
	"github.com/artheus/ohgossh/ssh"
	"github.com/artheus/ohgossh/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrorNoHostname = errors.New("No host provided")

	hostURL *url.URL
	command []string
)

var rootCommand = &cobra.Command{
	Use: "ohgossh [flags] destination [command]",
	Aliases: []string{
		"ssh",
		"gossh",
	},
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if flagsConfig.debugLogging {
			logrus.Infof("Activating debug logging")
			logrus.SetLevel(logrus.DebugLevel)
		}

		logrus.Debugf("args: %+v", args)

		if len(args) == 0 {
			return ErrorNoHostname
		}

		var hostString = args[0]

		if !strings.HasPrefix(hostString, "ssh://") {
			hostString = fmt.Sprintf("ssh://%s", hostString)
		}

		hostURL, err = url.Parse(hostString)
		if err != nil {
			return errors.Wrap(err, "unable to parse hostURL as url")
		}

		if len(args) > 1 {
			// Remove -- from command, if given as an argument before command
			if args[1] == "--" {
				args = args[1:]
			}

			command = args[1:]
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		defer utils.HandleError(&err)

		var currentUser *user.User
		currentUser, err = user.Current()
		utils.PanicOnError(err)

		if flagsConfig.config == "" {
			flagsConfig.config = "~/.ssh/ohgossh.yml"
		}

		if flagsConfig.knownHostsFile == "" {
			flagsConfig.knownHostsFile = "~/.ohgossh/known_hosts"
		}

		if strings.HasPrefix(flagsConfig.config, "~") {
			flagsConfig.config = filepath.Join(currentUser.HomeDir, flagsConfig.config[1:])
		}

		if strings.HasPrefix(flagsConfig.knownHostsFile, "~") {
			flagsConfig.knownHostsFile = filepath.Join(currentUser.HomeDir, flagsConfig.knownHostsFile[1:])
		}

		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		defer utils.HandleError(&err)

		cmd.Root().SilenceUsage = true
		cmd.Root().SilenceErrors = true

		config, err := config.LoadConfig(flagsConfig.config)
		utils.PanicOnError(err)

		config.KnownHostsFilename = flagsConfig.knownHostsFile

		host, err := host.Parse(hostURL, config)
		utils.PanicOnError(errors.Wrap(err, "failed to parse host from config"))

		host.URL = hostURL

		if hostURL.User != nil && hostURL.User.Username() != "" {
			host.User = hostURL.User.Username()
		}

		if hostURL.Port() != "" {
			var portnum uint64
			portnum, err = strconv.ParseUint(hostURL.Port(), 10, 16)
			utils.PanicOnError(errors.Wrapf(err, "unable to parse port number as uint16: %s", hostURL.Port()))
			host.Port = uint16(portnum)
		}

		utils.PanicOnError(
			errors.Wrap(
				ssh.Connect(host, command),
				"failed to connect to host",
			),
		)

		return err
	},
}

func Execute() error {
	return rootCommand.Execute()
}
