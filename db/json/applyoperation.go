package json

func (c *Client[Db]) ApplyOperation(f func(db *Db) *Db) {
	c.m.Lock()
	defer c.m.Unlock()
	c.db = *f(&c.db)
	c.autosave()
}
