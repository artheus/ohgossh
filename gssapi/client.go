package gssapi

type client struct {
}

func NewClient() *client {
	return &client{}
}

func (c *client) InitSecContext(target string, token []byte, isGSSDelegCreds bool) (outputToken []byte, needContinue bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) GetMIC(micFiled []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (c *client) DeleteSecContext() error {
	//TODO implement me
	panic("implement me")
}
