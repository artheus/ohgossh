package cmd

import (
	"github.com/artheus/ohgossh/version"
	"github.com/spf13/pflag"
)

var (
	flagsConfig = &FlagsConfig{}
)

type FlagsConfig struct {
	config         string
	knownHostsFile string
	verbose        int
}

func init() {
	// Note: Stop flag parsing after first non-flag argument
	rootCommand.Flags().SetInterspersed(false)

	rootCommand.Version = version.Version()

	flagSet := pflag.NewFlagSet("flags", pflag.ExitOnError)
	flagSet.CountVarP(&flagsConfig.verbose, "verbose", "v", "increased verbosity per flag added")

	flagSet.StringVarP(&flagsConfig.config, "config", "c", "~/.ssh/ohgossh.yml", "config file path")
	flagSet.StringVarP(&flagsConfig.knownHostsFile, "known-hosts", "k", "~/.ssh/known_hosts", "known_hosts file path")

	rootCommand.PersistentFlags().AddFlagSet(flagSet)
}
