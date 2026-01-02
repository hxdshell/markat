package core

import (
	"log"
)

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
