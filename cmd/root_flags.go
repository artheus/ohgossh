package cmd

import (
	"github.com/artheus/ohgossh/version"
)

var (
	flagsConfig = &FlagsConfig{}
)

type FlagsConfig struct {
	debugLogging   bool
	traceLogging   bool
	config         string
	knownHostsFile string
}

func init() {
	// Note: Stop flag parsing after first non-flag argument
	rootCommand.Flags().SetInterspersed(false)

	rootCommand.Version = version.Version()

	rootCommand.PersistentFlags().BoolVar(&flagsConfig.debugLogging, "debug", false, "enable debug logging")
	rootCommand.PersistentFlags().BoolVar(&flagsConfig.traceLogging, "trace", false, "enable trace logging")

	rootCommand.PersistentFlags().StringVarP(&flagsConfig.config, "config", "c", "~/.ssh/ohgossh.yml", "config file path")
	rootCommand.PersistentFlags().StringVarP(&flagsConfig.knownHostsFile, "known-hosts", "k", "~/.ssh/known_hosts", "known_hosts file path")
}
