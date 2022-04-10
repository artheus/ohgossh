package config

func DefaultHostParams() HostParams {
	return HostParams{
		Port:         22,
		IdentityFile: "~/.ssh/id_rsa",
		PreferredAuth: []string{
			"gssapi-with-mic",
			"publickey",
			"keyboard-interactive",
			"password",
		},
		Timeout:               0,
		InsecureIgnoreHostKey: false,
	}
}
