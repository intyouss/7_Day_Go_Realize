package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	uBytes    int64
	lList     *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		lList:     list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.lList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.lList.Back()
	if ele != nil {
		c.lList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.uBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.lList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.uBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.lList.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.uBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.uBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.lList.Len()
}
