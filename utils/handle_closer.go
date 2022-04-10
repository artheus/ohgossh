package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

func Close(c io.Closer) {
	_ = c.Close()
}

func CloseAndLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Errorf("%s", err)

		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			fmt.Printf("%+v\n", err)
		}
	}
}
