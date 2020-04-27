package collector

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"geochats/pkg/client"
	"geochats/pkg/collector/loaders"
	"geochats/pkg/storage"
	"regexp"
	"strconv"
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
		newGroup, err := c.loader.Export(g.Username)
		if err != nil {
			return fmt.Errorf("collector error: %v", err)
		}
		lat, long := checkCoordinates(newGroup.Description)
		newGroup.Latitude = lat
		newGroup.Longitude = long
		if err := c.store.AddGroup(newGroup); err != nil {
			return fmt.Errorf("collector error: %v", err)
		}
	}
	return nil
}

func checkCoordinates(desc string) (float64, float64) {
	var re = regexp.MustCompile(`(?m)https:\/\/miting.link\/#([0-9.]+),([0-9.]+)`)
	matches := re.FindAllStringSubmatch(desc, 1)
	if len(matches) == 0 {
		return 0, 0
	}
	lat, _ := strconv.ParseFloat(matches[0][1], 64)
	long, _ := strconv.ParseFloat(matches[0][2], 64)
	return lat, long

}
