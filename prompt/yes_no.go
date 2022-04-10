package prompt

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func YesNo(prompt string) (input bool, err error) {
	fmt.Printf("%s [Y/N]: ", prompt)

	var inputStr string

	input = false

	reader := bufio.NewReader(os.Stdin)
	inputStr, err = reader.ReadString('\n')
	fmt.Print("\n")

	inputStr = strings.TrimSpace(inputStr)

	if inputStr == "Y" {
		input = true
	}

	if err = errors.Wrap(err, "unable to prompt user for password"); err != nil {
		return
	}

	return
}
