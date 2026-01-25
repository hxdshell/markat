package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"markat/utils"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
)

type MessageText struct {
	Specifier   string
	ContentType string
	Encoding    string
}

type MessageAttachment struct {
	Specifier   string
	ContentType string
	Encoding    string
	FileName    string
	Size        string
}

// This structure contains all the leaf node metadata which helps to fetch them individually
type MessageStructure struct {
	MainType string
	Text     []MessageText

	// TextPlain and TextHtml are primary form of text rendering for client and contains index to Text, if not exists set to -1
	TextPlain   int
	TextHtml    int
	Attachments []MessageAttachment
}

func (c *Core) FetchMessageStructure(ctx context.Context, mb string, uid uint32) (*MessageStructure, error) {
	currentMb := c.ImapClient.GetCurrentMbName()

	if currentMb != mb {
		_, err := c.SelectMailBox(ctx, mb)
		if err != nil {
			return nil, err
		}
	}
	bs, err := c.ImapClient.FetchBodyStrucutre(ctx, uid)
	if err != nil {
		return nil, err
	}
	if bs == nil {
		return nil, errors.New("not found")
	}
	mainType := fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)
	ms := &MessageStructure{
		MainType:  mainType,
		TextPlain: -1,
		TextHtml:  -1,
	}
	walkBs(bs, ms, nil)

	return ms, nil
}

func (c *Core) FetchMessageText(ctx context.Context, mb string, uid uint32) ([]byte, error) {
	var b []byte

	ms, err := c.FetchMessageStructure(ctx, mb, uid)
	if err != nil {
		return b, err
	}
	if ms.TextPlain == -1 {
		// only deal with text/plain for now
		return b, errors.New("not found")
	}
	specifier := ms.Text[ms.TextPlain].Specifier
	msg, _, err := c.ImapClient.FetchMessage(ctx, specifier, uid)

	var r imap.Literal
	for _, literal := range msg.Body {
		r = literal // only one section
	}

	if strings.HasPrefix(ms.MainType, "text/") {
		entity, err := message.Read(r)
		if err != nil {
			return b, err
		}
		b, err = io.ReadAll(entity.Body)
		fmt.Println(entity.Body)
		if err != nil {
			return b, err
		}
		return b, nil
	}
	b, err = io.ReadAll(r)
	if err != nil {
		return b, err
	}

	return b, nil
}

func walkBs(bs *imap.BodyStructure, ms *MessageStructure, parts []int) {
	if bs == nil {
		return
	}

	mimeType := strings.ToLower(bs.MIMEType)
	subType := strings.ToLower(bs.MIMESubType)

	if mimeType != "multipart" {
		if mimeType == "text" {
			msgTxt := MessageText{}
			msgTxt.ContentType = fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)

			specifier := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(parts)), "."), "[]")

			msgTxt.Specifier = specifier
			msgTxt.Encoding = bs.Encoding

			ms.Text = append(ms.Text, msgTxt)

			switch subType {
			case "plain":
				ms.TextPlain = len(ms.Text) - 1
			case "html":
				ms.TextHtml = len(ms.Text) - 1
			}
		} else {
			msgAtchmnt := MessageAttachment{}
			msgAtchmnt.ContentType = fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)
			specifier := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(parts)), "."), "[]")
			msgAtchmnt.Specifier = specifier
			msgAtchmnt.Encoding = bs.Encoding
			filename, _ := bs.Filename() // ignoring the error for now
			msgAtchmnt.FileName = filename
			msgAtchmnt.Size = utils.HumanMessageSize(uint(bs.Size), false, 2)
			ms.Attachments = append(ms.Attachments, msgAtchmnt)
		}
	}

	for i, p := range bs.Parts {
		walkBs(p, ms, append(parts, i+1))
	}
}
