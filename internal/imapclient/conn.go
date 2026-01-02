package imapclient

import (
	"sync"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type ImapClient struct {
	conn *client.Client
	sync.RWMutex
	mb *imap.MailboxStatus
}

func StartTLS(addr string) (*ImapClient, error) {
	conn, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}
	return &ImapClient{conn: conn}, err
}

func (ic *ImapClient) Login(username string, password string) error {
	ic.Lock()
	err := ic.conn.Login(username, password)
	ic.Unlock()
	return err
}

func (ic *ImapClient) Logout() error {
	ic.Lock()
	err := ic.conn.Logout()
	ic.Unlock()
	return err
}

func (ic *ImapClient) Noop() {
	ic.Lock()
	ic.conn.Noop() // ignoring the error
	ic.Unlock()
}
