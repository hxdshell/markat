package core

import (
	"markat/internal/imapclient"
)

type Core struct {
	ImapClient *imapclient.ImapClient
}

func InitCore(client *imapclient.ImapClient) *Core {
	return &Core{
		ImapClient: client,
	}
}
