package imapclient

import (
	"context"
	"errors"

	"github.com/emersion/go-imap"
)

func (ic *ImapClient) List(ctx context.Context) ([]imap.MailboxInfo, error) {
	mbchan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.List("", "*", mbchan)
	}()

	var mailboxes []imap.MailboxInfo
	for {
		select {
		case mb, ok := <-mbchan:
			if !ok {
				return mailboxes, nil
			}
			mailboxes = append(mailboxes, *mb)
		case <-ctx.Done():
			return mailboxes, ctx.Err()
		}
	}
}

func (ic *ImapClient) Select(ctx context.Context, name string) (*imap.MailboxStatus, error) {
	done := make(chan bool, 1)
	var status *imap.MailboxStatus
	var err error
	ic.Lock()
	defer ic.Unlock()
	go func() {
		status, err = ic.conn.Select(name, false)
		ic.mb = status
		done <- true
	}()
	select {
	case <-done:
		return status, err
	case <-ctx.Done():
		return status, ctx.Err()
	}
}

func (ic *ImapClient) FetchEnvelopes(ctx context.Context, total uint32) ([]imap.Message, error) {
	var messages []imap.Message
	if ic.mb == nil {
		return messages, errors.New("mailbox is not selected")
	}
	from := uint32(1)
	to := ic.mb.Messages
	if ic.mb.Messages > total {
		from = ic.mb.Messages - total
	}
	seqset := &imap.SeqSet{}
	seqset.AddRange(from, to)
	items := []imap.FetchItem{imap.FetchUid, imap.FetchInternalDate, imap.FetchEnvelope, imap.FetchFlags, imap.FetchRFC822Size}

	msgchan := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.Fetch(seqset, items, msgchan)
	}()

	for {
		select {
		case msg, ok := <-msgchan:
			if !ok {
				return messages, nil
			}
			messages = append(messages, *msg)
		case <-ctx.Done():
			return messages, ctx.Err()
		}
	}
}
