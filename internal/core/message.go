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
	Specifier   string `json:"specifier"`
	ContentType string `json:"contentType"`
	Encoding    string `json:"encoding"`
	FileName    string `json:"fileName"`
	Size        string `json:"size"`
}

// This structure contains all the leaf node metadata which helps to fetch them individually
type MessageStructure struct {
	Uid      uint32
	Mb       string
	MainType string
	Text     []MessageText

	// TextPlain and TextHtml are primary form of text rendering for client and contains index to Text, if not exists set to -1
	TextPlain   int
	TextHtml    int
	Attachments []MessageAttachment
}

type MessageMeta struct {
	Uid          uint32              `json:"uid"`
	Mb           string              `json:"mb"`
	From         string              `json:"from"`
	To           string              `json:"to"`
	Attatchments []MessageAttachment `json:"attachments"`
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

func (c *Core) fetchMessageStructure(ctx context.Context, mb string, uid uint32) (*MessageStructure, error) {
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
		Uid:       uid,
		Mb:        mb,
		MainType:  mainType,
		TextPlain: -1,
		TextHtml:  -1,
	}
	walkBs(bs, ms, nil)

	return ms, nil
}

func (c *Core) FetchMeta(ctx context.Context, mb string, uid uint32) (*MessageMeta, error) {

	ms, err := c.fetchMessageStructure(ctx, mb, uid)
	if err != nil {
		return nil, err
	}
	msg, err := c.ImapClient.FetchHeader(ctx, uid)
	if err != nil {
		return nil, err
	}
	var r imap.Literal
	for _, literal := range msg.Body {
		r = literal // only one section
	}
	entity, err := message.Read(r)
	headers := entity.Header

	c.Lock()
	c.RecentMsgStructure = ms
	c.Unlock()

	meta := &MessageMeta{
		Uid:          uid,
		Mb:           mb,
		From:         headers.Get("From"),
		To:           headers.Get("To"),
		Attatchments: ms.Attachments,
	}
	return meta, nil
}
func (c *Core) msgStructureFromCache(ctx context.Context, mb string, uid uint32) (*MessageStructure, error) {
	if c.RecentMsgStructure != nil {
		if c.RecentMsgStructure.Uid == uid && c.RecentMsgStructure.Mb == mb {
			c.RLock()
			ms := c.RecentMsgStructure
			c.RUnlock()
			return ms, nil
		}
	}
	c.RecentMsgStructure = nil
	return c.fetchMessageStructure(ctx, mb, uid)
}

func (c *Core) FetchMessageText(ctx context.Context, mb string, uid uint32) ([]byte, error) {
	var b []byte

	ms, err := c.msgStructureFromCache(ctx, mb, uid)
	if err != nil {
		return b, err
	}

	if ms.TextPlain == -1 {
		// only deal with text/plain for now
		return b, errors.New("not found")
	}
	part := ms.Text[ms.TextPlain]
	specifier := part.Specifier
	msg, _, err := c.ImapClient.FetchMessage(ctx, specifier, uid)

	if err != nil {
		return b, err
	}

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
		if err != nil {
			return b, err
		}
		return b, nil
	}
	b, err = io.ReadAll(r)
	if err != nil {
		return b, err
	}

	if part.Encoding == "base64" {
		b, err = utils.DecodeBase64(string(b))
	}

	return b, nil
}
