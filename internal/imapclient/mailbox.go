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

func (ic *ImapClient) FetchEnvelopes(ctx context.Context, uids []uint32) ([]imap.Message, error) {
	var messages []imap.Message
	if ic.mb == nil {
		return messages, errors.New("mailbox is not selected")
	}

	seqset := &imap.SeqSet{}
	seqset.AddNum(uids...)
	items := []imap.FetchItem{imap.FetchUid, imap.FetchInternalDate, imap.FetchEnvelope, imap.FetchFlags, imap.FetchRFC822Size}

	msgchan := make(chan *imap.Message, 20)
	done := make(chan error, 1)
	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.UidFetch(seqset, items, msgchan)
	}()
	msgByUID := map[uint32]*imap.Message{}
	select {
	case err := <-done:
		if err != nil {
			return messages, err
		}
		for msg := range msgchan {
			msgByUID[msg.Uid] = msg
		}
		for _, uid := range uids {
			msg := msgByUID[uid]
			messages = append(messages, *msg)
		}
		return messages, nil
	case <-ctx.Done():
		return messages, ctx.Err()
	}
}

func (ic *ImapClient) FetchAllUids(ctx context.Context) ([]imap.Message, error) {
	var messages []imap.Message
	if ic.mb == nil {
		return messages, errors.New("mailbox is not selected")
	}
	seqset := &imap.SeqSet{}
	seqset.AddRange(1, ic.mb.Messages)
	items := []imap.FetchItem{imap.FetchUid, imap.FetchInternalDate}

	msgchan := make(chan *imap.Message, 20)
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
