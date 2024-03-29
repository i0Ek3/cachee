package singleflight

import "sync"

type call struct {
	wg sync.WaitGroup
	val any
	err error
}

type Group struct {
	mu sync.Mutex
	m map[string]*call
}

// Do only do once
func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.m, key)

	return c.val, c.err
}