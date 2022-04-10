package prompt

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func PlainText(prompt string) (input string, err error) {
	fmt.Printf("%s: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err = reader.ReadString('\n')
	if err = errors.Wrap(err, "unable to prompt user for password"); err != nil {
		return "", err
	}
	fmt.Print("\n")

	input = strings.TrimSpace(input)

	return
}
