package collector

import (
	"github.com/sirupsen/logrus"
	"geochats/pkg/client"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
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
	return nil
}
