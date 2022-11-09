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
		if flagsConfig.verbose == 1 {
			logrus.Infof("Debug logging enabled")
			logrus.SetLevel(logrus.DebugLevel)
		} else if flagsConfig.verbose > 1 {
			logrus.Infof("Trace logging enabled")
			logrus.SetLevel(logrus.TraceLevel)
		}

		logrus.Tracef("Started with arguments: %+v", args)

		if len(args) == 0 {
			return ErrorNoHostname
		}

		var hostString string
		hostString, command = args[0], args[1:]

		logrus.Tracef("host: %s", hostString)
		logrus.Tracef("command: %+v", command)

		if !strings.HasPrefix(hostString, "ssh://") {
			hostString = fmt.Sprintf("ssh://%s", hostString)
		}

		hostURL, err = url.Parse(hostString)
		if err != nil {
			return errors.Wrap(err, "unable to parse hostURL as url")
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
		cmd.Root().SilenceUsage = true
		cmd.Root().SilenceErrors = true

		config, err := config.LoadConfig(flagsConfig.config)
		if err != nil {
			logrus.Errorf("Failed to load config: %+v", errors.WithStack(err))
		}

		config.KnownHostsFilename = flagsConfig.knownHostsFile

		host, err := host.Parse(hostURL, config)
		if err != nil {
			logrus.Errorf("failed to parse host from config: %+v", errors.WithStack(err))
		}

		logrus.Tracef("host url: %+v", hostURL)

		host.URL = hostURL

		if hostURL.User != nil && hostURL.User.Username() != "" {
			host.User = hostURL.User.Username()
		}

		if hostURL.Port() != "" {
			var portnum uint64
			if portnum, err = strconv.ParseUint(hostURL.Port(), 10, 16); err != nil {
				logrus.Errorf("unable to parse port number as uint16: %s: %+v", hostURL.Port(), errors.WithStack(err))
			}
			host.Port = uint16(portnum)
		}

		if err = ssh.Connect(host, command); err != nil {
			logrus.Errorf("connection failed: %+v", err)
		}

		return err
	},
}

func Execute() error {
	return rootCommand.Execute()
}
