package imapclient

import (
	"github.com/emersion/go-imap"
)

func (ic *ImapClient) List() ([]imap.MailboxInfo, error) {
	mbchan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.List("", "*", mbchan)
	}()

	var mailboxes []imap.MailboxInfo
	for mb := range mbchan {
		mailboxes = append(mailboxes, *mb)
	}
	err := <-done
	return mailboxes, err
}

func (ic *ImapClient) Select(name string) (*imap.MailboxStatus, error) {
	ic.Lock()
	status, err := ic.conn.Select(name, false)
	ic.mb = status
	ic.Unlock()
	return status, err
}
