package mail_client

import (
	"github.com/emersion/go-sasl"
)


type xoauth2Client struct {
	Username string
	Token    string
}


func NewXOAUTH2Client(username, token string) sasl.Client {
	return &xoauth2Client{
		Username: username,
		Token:    token,
	}
}


func (c *xoauth2Client) Start() (mech string, resp []byte, err error) {
	mech = "XOAUTH2"
	resp = []byte("user=" + c.Username + "\x01auth=Bearer " + c.Token + "\x01\x01")
	return
}


func (c *xoauth2Client) Next(challenge []byte) ([]byte, error) {
	return nil, nil
}
