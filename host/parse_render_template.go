package host

import (
	"bytes"
	"fmt"
	"github.com/artheus/ohgossh/prompt"
	"github.com/artheus/ohgossh/utils"
	regexp "github.com/gijsbers/go-pcre"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os/user"
	"reflect"
	"text/template"
)

func renderTemplate(hostname, templateStr string, captureGroups *regexp.Matcher) (renderedHostname string, err error) {
	defer utils.HandleError(&err)

	logrus.Debugf("rendering template: %s", templateStr)

	temp := template.New("hostname")

	logrus.Trace("getting current user")
	currentUser, err := user.Current()
	utils.PanicOnError(errors.Wrap(err, "not able to get current user"))

	logrus.Trace("creating template context")
	var tempCtx = &TemplateContext{
		ShellUser: currentUser.Username,
		User:      currentUser.Username,
		Host:      hostname,
	}

	logrus.Debugf("Template context: %+v", tempCtx)

	logrus.Trace("registering functions for template engine")
	var tempFuncs = map[string]any{
		"randomChoice": func(choices ...any) any {
			return choices[rand.Intn(len(choices))]
		},
		"askpass": promptForInput(true),
		"prompt":  promptForInput(false),
		"default": defaultFunc(),
	}

	if captureGroups != nil {
		for i := 1; i <= captureGroups.Groups(); i++ {
			tempFuncs[fmt.Sprintf("c%d", i)] = func(group int) func() string {
				return func() string {
					return captureGroups.GroupString(group)
				}
			}(i)
		}
	}

	temp.Funcs(tempFuncs)

	logrus.Trace("parsing replace template")
	_, err = temp.Parse(templateStr)
	utils.PanicOnError(errors.Wrap(err, "failed to parse hosts replace template"))

	logrus.Trace("rendering hostname from template")
	buf := bytes.NewBuffer([]byte{})

	utils.PanicOnError(
		temp.Execute(buf, &tempCtx),
	)

	renderedHostname = buf.String()

	logrus.Debugf("template rendered as: %s", renderedHostname)

	return renderedHostname, nil
}

func promptForInput(secret bool) func(string) string {
	return func(promptMsg string) (input string) {
		defer utils.LogErrors()

		var err error

		logrus.Debug("prompting user for input")

		if secret {
			input, err = prompt.Secret(promptMsg)
		} else {
			input, err = prompt.PlainText(promptMsg)
		}

		utils.PanicOnError(err)

		return input
	}
}

func defaultFunc() func(interface{}, interface{}) interface{} {
	return func(arg, value interface{}) interface{} {
		defer utils.LogErrors()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				return arg
			}
		case reflect.Bool:
			if !v.Bool() {
				return arg
			}
		default:
			return value
		}

		return value
	}
}
