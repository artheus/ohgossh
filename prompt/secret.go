package prompt

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/term"
	"strings"
	"syscall"
)

func Secret(prompt string) (secret string, err error) {
	fmt.Printf("%s: ", prompt)

	passBytes, err := term.ReadPassword(syscall.Stdin)
	fmt.Print("\n")

	if err = errors.Wrap(err, "unable to prompt user for password"); err != nil {
		return "", err
	}

	secret = strings.TrimSpace(string(passBytes))

	return secret, err
}
