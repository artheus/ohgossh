[![Go](https://github.com/artheus/ohgossh/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/artheus/ohgossh/actions/workflows/go.yml)

# ohgossh

SSH alternative command line client written in Go

## Why use ohgossh?

While configuring my ssh_config for the Nth time, I just figured
"Why not make an alternative client, which will allow Regexp patterns
for matching hostnames, and using a template for rendering hostname and
other options to use!?"

So I did. I made this project in Go. Using the `golang.org/x/crypto`
library, which provides most of what is needed for actually connecting
to an SSH server, and `golang.org/x/term` for some terminal stuff.

## What is implemented?

- [x] Simple ssh command line interface
- [x] Regexp pattern matching of hosts
- [x] Go template rendering of hostname and options
- [x] Piping the ssh connection to terminal, and making it actually usable
- [ ] GSSAPI authentication
- [ ] HTTP proxying of ssh connection

## Config format

ohgossh uses the yaml format for the config file. The yaml config allows
users to add their own PERL compatible regexp patterns and using Go templates
for rendering hostnames and options.

Here is an example config from `config/examples/regexp.yml`:

```yaml
---

defaults:
  user: "{{ default .User .ShellUser }}"
  timeout: "10s"
  identityFile: "~/.ssh/id_rsa{{if ne .User .ShellUser}}.{{.User}}{{end}}"
  preferredAuthentications:
    - gssapi-with-mic
    - publickey
    - password

hosts:
  - pattern: "jump-(uk|us)\\.(dev|stage)"
    replace: "jumphost-0{{if eq (capture 1) `uk`}}{{ random.Number 1 2 }}{{else}}{{ random.Number 3 4 }}{{end}}.{{capture 2}}.example.domain"
    user: admin

  - pattern: "jump-(uk|us)\\.prod"
    replace: "jumpprod-{{if eq (capture 1) `uk`}}{{ random.Number 1 2 }}{{else}}{{ random.Number 3 4 }}{{end}}.example.domain"
    user: "prod-admin"
    gssApi:
      enabled: true
      delegateCredentials: true

  - pattern: "k8s-((master|worker)[0-9]*)-(uk|us)\\.(dev|stage|prod)"
    replace: "{{capture 1}}-{{capture 2}}-{{capture 4}}.{{capture 5}}"
    user: "core"
    identityFile: "~/.ssh/id_rsa.ansible"

  - pattern: "((?!jump).*?-(uk|us))\\.(dev|stage|prod)"
    replace: "{{capture 1}}"
    jumpHost: "jump-{{capture 2}}.{{capture 3}}"

  - name: "github.com"
    aliases:
      - github
    httpProxy:
      host: proxy.localhost
      port: 8080
```
