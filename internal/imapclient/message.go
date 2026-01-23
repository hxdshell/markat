package imapclient

import (
	"context"

	"github.com/emersion/go-imap"
)

func (ic *ImapClient) FetchMessage(ctx context.Context, uid uint32) (*imap.Message, error) {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)
	items := []imap.FetchItem{imap.FetchRFC822}
	msgchan := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- ic.conn.UidFetch(seqset, items, msgchan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-done:
		return nil, err
	case msg := <-msgchan:
		return msg, nil
	}
}
