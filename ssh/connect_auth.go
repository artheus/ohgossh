package ssh

import (
	"crypto/x509"
	"fmt"
	"github.com/artheus/ohgossh/host"
	"github.com/artheus/ohgossh/prompt"
	"github.com/bodgit/sshkrb5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func authMethodFor(auth string, host *host.Host) (authMethod ssh.AuthMethod, err error) {
	switch auth {
	case "password":
		var pass string
		pass, err = prompt.Secret(fmt.Sprintf("Enter password for %s@%s", host.User, host.Name))
		return ssh.Password(pass), errors.Wrap(err, "password prompt failed")
	case "publickey":
		var pemBytes []byte
		var signer ssh.Signer

		pemBytes, err = ioutil.ReadFile(host.IdentityFile)
		if err != nil {
			return nil, errors.Wrapf(err, "could not read identityfile %s", host.IdentityFile)
		}

		signer, err = ssh.ParsePrivateKey(pemBytes)
		if _, ok := err.(*ssh.PassphraseMissingError); ok {
			logrus.Warnf("failed to parse private key: %+v", err)
			signer, err = privateKeyPassphraseLoop(host, pemBytes, 1)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get passphrase for private key")
			}
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to parse private key")
		}

		return ssh.PublicKeys(signer), err
	case "gssapi-with-mic":
		var gaClient *sshkrb5.Client

		if gaClient, err = sshkrb5.NewClient(); err != nil {
			return nil, errors.Wrap(err, "GSSAPI client failed to initialize")
		}

		return ssh.GSSAPIWithMICAuthMethod(gaClient, host.Name), nil
	case "keyboard-interactive":
		return ssh.KeyboardInteractive(keyboardChallenge), nil
	}

	return nil, errors.Errorf("unsupported auth method: %s", auth)
}

func keyboardChallenge(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	fmt.Printf("Keyboard challenge authentication for %s", user)
	fmt.Println(instruction)
	fmt.Println("\nQuestions:")

	var answer string

	for i, question := range questions {
		fmt.Printf("%d. %s\n", i, question)
		if echos[i] {
			answer, err = prompt.PlainText("Answer")
		} else {
			answer, err = prompt.Secret("Secret answer")
		}

		if err != nil {
			return []string{}, errors.Wrap(err, "keyboard challenge failed")
		}

		answers = append(answers, answer)
	}

	return
}

const maxPassphraseRetries = 3

func privateKeyPassphraseLoop(host *host.Host, pemBytes []byte, tryNum int) (signer ssh.Signer, err error) {
	var passphrase string

	passphrase, err = prompt.Secret(fmt.Sprintf("Passphrase for %s", host.IdentityFile))

	if err != nil {
		return nil, errors.Wrap(err, "failed to prompt user for passphrase")
	}

	if signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(passphrase)); err != nil {
		if err == x509.IncorrectPasswordError {
			if tryNum >= maxPassphraseRetries {
				return nil, errors.New("tried passphrase for private key 3 times without luck..")
			}

			fmt.Println("Incorrect passphrase, please try again")
			return privateKeyPassphraseLoop(host, pemBytes, tryNum+1)
		} else {
			return nil, err
		}
	}

	return
}
