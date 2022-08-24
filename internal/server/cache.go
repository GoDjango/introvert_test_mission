package server

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	updatePeriod time.Duration
	entities     []Entity
	mdb          *MongoDB
}

func NewCache(updatePeriod time.Duration, mdb *MongoDB) *Cache {
	c := Cache{
		updatePeriod: updatePeriod,
		entities:     make([]Entity, 0),
		mdb:          mdb,
	}
	go c.Run()
	return &c
}

func (c *Cache) sync() {
	entities, err := c.mdb.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	c.Lock()
	defer c.Unlock()

	c.entities = entities
}

func (c *Cache) Run() {
	for {
		c.sync()
		time.Sleep(c.updatePeriod)
	}
}

func (c *Cache) GetAll() []Entity {
	c.RLock()
	defer c.RUnlock()
	return c.entities
}
