package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"markat/utils"
	"sort"
)

type MailBoxInfo struct {
	Messages uint32 `json:"messages"`
	Name     string `json:"name"`
	ReadOnly bool   `json:"readOnly"`
	Recent   uint32 `json:"recent"`
	Unseen   uint32 `json:"unseen"`
}

type EnvelopeDisplay struct {
	Uid      uint32   `json:"uid"`
	From     []string `json:"from"`
	FromName []string `json:"fromName"`
	To       []string `json:"to"`
	ToName   []string `json:"toName"`
	Sender   []string `json:"sender"`
	Date     string   `json:"date"`
	Subject  string   `json:"subject"`
	Size     string   `json:"size"`
	Flags    []string `json:"flags"`
}

type EnvelopeResponse struct {
	Page      int               `json:"page"`
	Start     int               `json:"start"`
	End       int               `json:"end"`
	Total     int               `json:"total"`
	Envelopes []EnvelopeDisplay `json:"envelopes"`
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

func (c *Core) FetchEnvelopes(ctx context.Context, page int, pageSize int) (*EnvelopeResponse, error) {
	envelopes := []EnvelopeDisplay{}
	response := &EnvelopeResponse{}
	response.Envelopes = envelopes
	messages, err := c.ImapClient.FetchAllUids(ctx)

	if err != nil {
		return response, err
	} else if len(messages) == 0 {
		return response, err
	}

	sort.Slice(messages, func(i, j int) bool {
		val := messages[i].InternalDate.Compare(messages[j].InternalDate)

		switch val {
		case 1:
			return true
		case -1:
			return false
		}

		return messages[i].Uid > messages[j].Uid
	})

	var uids []uint32
	for _, msg := range messages {
		uids = append(uids, msg.Uid)
	}
	start := pageSize * (page - 1)

	if start >= len(uids) || start < 0 {
		return response, errors.New("404")
	}

	end := min(start+pageSize, len(uids))

	slicedUids := uids[start:end]
	messages, err = c.ImapClient.FetchEnvelopes(ctx, slicedUids)
	if err != nil {
		return response, err
	}
	response.Page = page
	response.Start = start + 1
	response.End = end
	response.Total = len(uids)
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

		readableDate := msg.Envelope.Date.Format("02 Jan 2006 15:04")
		envlp := EnvelopeDisplay{
			Uid:      msg.Uid,
			Subject:  msg.Envelope.Subject,
			From:     fromList,
			FromName: fromNameList,
			To:       toList,
			ToName:   toNameList,
			Sender:   senderList,
			Date:     readableDate,
			Size:     utils.HumanMessageSize(uint(msg.Size), true, 2),
			Flags:    msg.Flags,
		}
		envelopes = append(envelopes, envlp)
	}

	log.Println("ENVELOPES", len(envelopes))
	response.Envelopes = envelopes
	return response, err
}
