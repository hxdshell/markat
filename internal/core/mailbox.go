package core

import (
	"log"
)

type MailBoxInfo struct {
	Messages uint32 `json:"messages"`
	Name     string `json:"name"`
	ReadOnly bool   `json:"readOnly"`
	Recent   uint32 `json:"recent"`
	Unseen   uint32 `json:"unseen"`
}

func (c *Core) MailBoxList() ([]string, error) {
	mailboxes, err := c.ImapClient.List()
	if err == nil {
		log.Println("LIST MAILBOX", len(mailboxes))
	}

	var mbnames []string
	for _, mb := range mailboxes {
		mbnames = append(mbnames, mb.Name)
	}
	return mbnames, err
}

func (c *Core) SelectMailBox(name string) (*MailBoxInfo, error) {
	var mbInfo *MailBoxInfo
	mb, err := c.ImapClient.Select(name)
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
