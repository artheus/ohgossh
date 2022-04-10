package utils

import (
	"github.com/sirupsen/logrus"
)

func HandleError(err *error) {
	if v := recover(); v != nil {
		if recErr, ok := v.(error); ok && recErr != nil {
			*err = recErr
		} else {
			panic(v)
		}
	}
}

func LogErrors() {
	if v := recover(); v != nil {
		if err, ok := v.(error); ok && err != nil {
			logrus.Errorf("runtime error: %v", err)
			panic(err)
		} else {
			panic(v)
		}
	}
}
