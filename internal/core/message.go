package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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

type MessageAttachmentHeader struct {
	Specifier   string
	Type        string
	Disposition string
	Encoding    string
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
	Uid         uint32              `json:"uid"`
	Mb          string              `json:"mb"`
	From        string              `json:"from"`
	To          string              `json:"to"`
	Subject     string              `json:"subject"`
	Date        string              `json:"date"`
	Attachments []MessageAttachment `json:"attachments"`
}

func (c *Core) walkBs(bs *imap.BodyStructure, ms *MessageStructure, parts []int) {
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

			contentDisposition := bs.Disposition
			contentType := fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)
			contentLength := bs.Size

			msgAtchmnt.ContentType = contentType
			specifier := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(parts)), "."), "[]")
			msgAtchmnt.Specifier = specifier
			msgAtchmnt.Encoding = bs.Encoding
			filename, _ := bs.Filename() // ignoring the error for now
			msgAtchmnt.FileName = filename
			msgAtchmnt.Size = utils.HumanMessageSize(uint(contentLength), false, 2)
			ms.Attachments = append(ms.Attachments, msgAtchmnt)

			// Cache the headers for the purposes of serving the file
			attchmntHeader := MessageAttachmentHeader{}

			attchmntHeader.Specifier = specifier

			attchmntHeader.Disposition = fmt.Sprintf("%s; filename=\"%s\"", contentDisposition, filename)

			attchmntHeader.Type = contentType
			attchmntHeader.Encoding = bs.Encoding

			c.Lock()
			c.CurrentAttachments = append(c.CurrentAttachments, attchmntHeader)
			c.Unlock()

		}
	}

	for i, p := range bs.Parts {
		c.walkBs(p, ms, append(parts, i+1))
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
	// Clear attachments first
	c.Lock()
	c.CurrentAttachments = nil
	c.Unlock()

	// and then walk bs
	c.walkBs(bs, ms, nil)

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

	readableDate := msg.Envelope.Date.Format("Mon, 02 Jan 2006 15:04")
	meta := &MessageMeta{
		Uid:         uid,
		Mb:          mb,
		From:        headers.Get("From"),
		To:          headers.Get("To"),
		Subject:     msg.Envelope.Subject,
		Date:        readableDate,
		Attachments: ms.Attachments,
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
	if msg == nil {
		return b, errors.New("not found")
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
		b, err = utils.DecodeBase64(b)
	}
	log.Println("MESSAGE:", mb, uid)
	return b, nil
}

func (c *Core) FetchAttachment(ctx context.Context, mb string, uid uint32, specifier string) (MessageAttachmentHeader, []byte, error) {
	var b []byte
	var h MessageAttachmentHeader

	currentMb := c.ImapClient.GetCurrentMbName()

	if currentMb != mb {
		_, err := c.SelectMailBox(ctx, mb)
		if err != nil {
			return h, b, err
		}
	}

	c.RLock()
	ms := c.RecentMsgStructure
	if ms != nil {
		if ms.Uid != uid {
			c.RUnlock()
			return h, b, errors.New("not found")
		}
	}
	headers := c.CurrentAttachments
	c.RUnlock()
	if headers == nil {
		return h, b, errors.New("not found")
	}

	found := false
	for _, h = range headers {
		if h.Specifier == specifier {
			found = true
			break
		}
	}
	if !found {
		return h, b, errors.New("not found")
	}

	msg, _, err := c.ImapClient.FetchMessage(ctx, specifier, uid)
	if err != nil {
		return h, b, err
	}
	if msg == nil {
		return h, b, errors.New("not found")
	}

	var r imap.Literal
	for _, literal := range msg.Body {
		r = literal // only one section
	}
	b, err = io.ReadAll(r)
	if err != nil {
		return h, b, err
	}
	if h.Encoding == "base64" {
		b, err = utils.DecodeBase64(b)
		if err != nil {
			return h, b, err
		}
	}
	log.Println("ATTACHMENT:", mb, uid, specifier)
	return h, b, nil
}
