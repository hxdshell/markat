package core

import (
	"context"
	"fmt"
	"log"
	"markat/utils"
	"time"
)

type MailBoxInfo struct {
	Messages uint32 `json:"messages"`
	Name     string `json:"name"`
	ReadOnly bool   `json:"readOnly"`
	Recent   uint32 `json:"recent"`
	Unseen   uint32 `json:"unseen"`
}

type EnvelopeDisplay struct {
	Uid          uint32    `json:"uid"`
	InternalDate time.Time `json:"internalDate"`
	From         []string  `json:"from"`
	FromName     []string  `json:"fromName"`
	To           []string  `json:"to"`
	ToName       []string  `json:"toName"`
	Sender       []string  `json:"sender"`
	Date         time.Time `json:"date"`
	Subject      string    `json:"subject"`
	Size         string    `json:"size"`
	Flags        []string  `json:"flags"`
}

func (c *Core) MailBoxList(ctx context.Context) ([]string, error) {
	mailboxes, err := c.ImapClient.List(ctx)
	if err == nil {
		log.Println("LIST MAILBOX", len(mailboxes))
	}

	var mbnames []string
	for _, mb := range mailboxes {
		mbnames = append(mbnames, mb.Name)
	}
	return mbnames, err
}

func (c *Core) SelectMailBox(ctx context.Context, name string) (*MailBoxInfo, error) {
	var mbInfo *MailBoxInfo
	mb, err := c.ImapClient.Select(ctx, name)
	if err == nil {
		log.Println("SELECT MAILBOX:", name)
		mbInfo = &MailBoxInfo{
			Messages: mb.Messages,
			Name:     mb.Name,
			ReadOnly: mb.ReadOnly,
			Recent:   mb.Recent,
			Unseen:   mb.Unseen,
		}
	}
	return mbInfo, err
}

func (c *Core) FetchEnvelopes(ctx context.Context) ([]EnvelopeDisplay, error) {
	var envelopes []EnvelopeDisplay
	messages, err := c.ImapClient.FetchEnvelopes(ctx, uint32(10))

	if err != nil {
		return envelopes, err
	}

	for _, msg := range messages {

		var fromList []string
		var fromNameList []string
		for _, address := range msg.Envelope.From {
			from := fmt.Sprintf("%s@%s", address.MailboxName, address.HostName)
			fromList = append(fromList, from)
			fromNameList = append(fromNameList, address.PersonalName)
		}

		var toList []string
		var toNameList []string
		for _, address := range msg.Envelope.To {
			to := fmt.Sprintf("%s@%s", address.MailboxName, address.HostName)
			toList = append(toList, to)
			toNameList = append(toNameList, address.PersonalName)
		}

		var senderList []string
		for _, address := range msg.Envelope.Sender {
			sender := fmt.Sprintf("%s@%s", address.MailboxName, address.HostName)

			senderList = append(senderList, sender)
		}

		envlp := EnvelopeDisplay{
			Uid:          msg.Uid,
			InternalDate: msg.InternalDate,
			Subject:      msg.Envelope.Subject,
			From:         fromList,
			FromName:     fromList,
			To:           toList,
			ToName:       toNameList,
			Sender:       senderList,
			Date:         msg.Envelope.Date,
			Size:         utils.HumanMessageSize(uint(msg.Size), true, 2),
			Flags:        msg.Flags,
		}
		envelopes = append(envelopes, envlp)
	}
	if envelopes == nil {
		envelopes = []EnvelopeDisplay{}
	} else {
		log.Println("ENVELOPES", len(envelopes))
	}
	return envelopes, err
}
