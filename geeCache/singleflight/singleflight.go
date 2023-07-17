package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type Group struct {
	mutex sync.Mutex
	m     map[string]*call
}

func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mutex.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mutex.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mutex.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mutex.Lock()
	delete(g.m, key)
	g.mutex.Unlock()

	return c.val, c.err
}
