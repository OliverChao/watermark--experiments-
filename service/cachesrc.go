package service

import (
	"WaterMasking/model"
	"sync"
)

var SourceCache = &cached{
	cacheFlag: false,
	mutex:     &sync.Mutex{},
}
var MarkedCache = &cached{
	cacheFlag: false,
	mutex:     &sync.Mutex{},
}

type cached struct {
	cacheFlag bool
	mutex     *sync.Mutex
	data      []*model.Student
}

func (c *cached) EnCached(d []*model.Student) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheFlag = true
	c.data = d
}

func (c *cached) UnCached() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheFlag = false
	c.data = []*model.Student{}
}

func (c *cached) IsCached() bool {
	return c.cacheFlag
}
func (c *cached) GetSourceData() []*model.Student {
	return c.data
}
