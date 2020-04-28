package collector

import (
	"fmt"
	"geochats/pkg/client"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"github.com/sirupsen/logrus"
)

type Collector struct {
	cl client.AbstractClient
	loader *loaders.ChannelInfoLoader
	store storage.Storage
	logger *logrus.Logger
}

func New(cl client.AbstractClient, loader *loaders.ChannelInfoLoader, store storage.Storage, logger *logrus.Logger) *Collector {
	return &Collector{
		cl: cl,
		loader: loader,
		store: store,
		logger: logger,
	}
}

func (c *Collector) UpdateGroups() error {
	groups, err := c.store.ListGroups()
	if err != nil {
		return fmt.Errorf("can't list groups from collector")
	}
	for _, g := range groups {
		chat, err := c.cl.GetChat(g.ChatID)
		if err != nil {
			c.logger.Errorf("can't get chat from tg: %v", err)
			continue
		}
		newGroup, err := c.loader.Export(chat, true)
		if err != nil {
			c.logger.Errorf("collector error: %v", err)
			continue
		}
		if _, err := c.store.UpdateGroup(newGroup); err != nil {
			c.logger.Errorf("collector error: %v", err)
			continue
		}
	}
	return nil
}
