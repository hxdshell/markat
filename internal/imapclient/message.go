package imapclient

import (
	"context"
	"errors"

	"github.com/emersion/go-imap"
)

func (ic *ImapClient) FetchBodyStrucutre(ctx context.Context, uid uint32) (*imap.BodyStructure, error) {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)
	items := []imap.FetchItem{imap.FetchBodyStructure}
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
		if msg == nil {
			return nil, errors.New("not found")
		}
		return msg.BodyStructure, nil
	}
}

func (ic *ImapClient) FetchHeader(ctx context.Context, uid uint32) (*imap.Message, error) {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)
	items := []imap.FetchItem{imap.FetchRFC822Header, imap.FetchEnvelope}
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
		if msg == nil {
			return nil, errors.New("not found")
		}
		return msg, nil
	}
}

func (ic *ImapClient) FetchMessage(ctx context.Context, specifier string, uid uint32) (*imap.Message, *imap.BodySectionName, error) {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)

	bodySection := &imap.BodySectionName{
		Peek: true,
		BodyPartName: imap.BodyPartName{
			Specifier: imap.PartSpecifier(specifier),
		},
	}
	items := []imap.FetchItem{bodySection.FetchItem()}
	msgchan := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.UidFetch(seqset, items, msgchan)
	}()

	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case err := <-done:
		return nil, nil, err
	case msg := <-msgchan:
		return msg, bodySection, nil
	}
}

func (ic *ImapClient) FetchMime(ctx context.Context, specifier string, uid uint32) (*imap.Message, *imap.BodySectionName, error) {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)

	bodySection := &imap.BodySectionName{
		Peek: true,
		BodyPartName: imap.BodyPartName{
			Specifier: imap.PartSpecifier(specifier + "." + imap.MIMESpecifier),
		},
	}
	items := []imap.FetchItem{bodySection.FetchItem()}
	msgchan := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	ic.Lock()
	defer ic.Unlock()
	go func() {
		done <- ic.conn.UidFetch(seqset, items, msgchan)
	}()

	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case err := <-done:
		return nil, nil, err
	case msg := <-msgchan:
		return msg, bodySection, nil
	}
}
func (ic *ImapClient) MarkSeenUnseen(ctx context.Context, uid uint32, seen bool) error {
	seqset := &imap.SeqSet{}
	seqset.AddNum(uid)

	var flagOp imap.FlagsOp = imap.RemoveFlags
	if seen {
		flagOp = imap.AddFlags
	}
	// this will be empty since it's a silent operation
	msgchan := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- ic.conn.UidStore(seqset, imap.FormatFlagsOp(flagOp, true), []any{imap.SeenFlag}, msgchan)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
