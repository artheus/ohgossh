package host

import "io"

type passwordAuth struct {
}

func (p2 *passwordAuth) auth(session []byte, user string, p interface{}, rand io.Reader) (interface{}, []string, error) {
	//TODO implement me
	panic("implement me")
}

func (p2 *passwordAuth) method() string {
	//TODO implement me
	panic("implement me")
}
