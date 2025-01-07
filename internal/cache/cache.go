package cache

import "log/slog"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type Element struct {
	Key   Key
	Value any
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem

	log *slog.Logger
}

func New(capacity int, log *slog.Logger) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		log:      log,
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = &Element{
			Key:   key,
			Value: value,
		}
		c.queue.MoveToFront(item)

		return true
	}

	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		element, ok := lastItem.Value.(*Element)
		if !ok {
			c.log.Error("invalid element type")
			return false
		}
		delete(c.items, element.Key)
	}

	c.items[key] = c.queue.PushFront(
		&Element{
			Key:   key,
			Value: value,
		},
	)

	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)

		element, ok := item.Value.(*Element)
		if !ok {
			c.log.Error("invalid element type")
			return nil, false
		}

		return element.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
