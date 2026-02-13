package core

import (
	"markat/internal/imapclient"
	"sync"
)

type Core struct {
	sync.RWMutex
	ImapClient         *imapclient.ImapClient
	RecentMsgStructure *MessageStructure
	CurrentAttachments []MessageAttachmentHeader
}

func InitCore(client *imapclient.ImapClient) *Core {
	return &Core{
		ImapClient:         client,
		CurrentAttachments: nil,
	}
}
