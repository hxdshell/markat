package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
)

type MessagePart struct {
	Header message.Header
	Body   []byte
}

func (c *Core) FetchMessage(ctx context.Context, mailbox string, uid uint32) ([]MessagePart, error) {
	currentMb := c.ImapClient.GetCurrentMbName()
	parts := []MessagePart{}

	if currentMb != mailbox {
		_, err := c.SelectMailBox(ctx, mailbox)
		if err != nil {
			return parts, err
		}
	}

	msg, err := c.ImapClient.FetchMessage(ctx, uid)
	if err != nil {
		return parts, err
	}
	if msg == nil {
		return parts, errors.New("not found")
	}
	bodySection := &imap.BodySectionName{Peek: true}
	r := msg.GetBody(bodySection)

	entity, err := message.Read(r)
	err = iterateMsgBody(&parts, entity)
	if err != nil {
		log.Println(err)
		return parts, err
	}
	return parts, nil
}

func iterateMsgBody(parts *[]MessagePart, entity *message.Entity) error {
	if mr := entity.MultipartReader(); mr != nil {
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			if err := iterateMsgBody(parts, p); err != nil {
				return err
			}
		}
		return nil
	} else {
		body := entity.Body
		fmt.Println(entity.Header.Get("Content-Type"))
		if strings.HasPrefix(entity.Header.Get("Content-Type"), "text/") {
			b, err := io.ReadAll(body)
			if err != nil {
				return err
			}
			part := &MessagePart{
				Header: entity.Header,
				Body:   b,
			}
			*parts = append(*parts, *part)
		}
		return nil
	}
}
