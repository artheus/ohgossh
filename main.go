package main

import (
	"fmt"
	"github.com/artheus/ohgossh/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			fmt.Printf("%+v\n", err)
		}

		os.Exit(1)
	}
}
