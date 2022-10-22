package version

import _ "embed"

//go:embed version.txt
var _version string

func Version() string {
	return _version
}
