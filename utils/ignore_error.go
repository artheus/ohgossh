package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

func IgnoreErr(_ error) {
	// do nothing
}

func IgnoreButLogErr(err error) {
	logrus.Errorf("%s", err)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		fmt.Printf("%+v\n", err)
	}
}

func IgnoreErrIOCopy(dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
}

func IgnoreButLogErrIOCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)

	logrus.Errorf("%s", err)

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		fmt.Printf("%+v\n", err)
	}
}
