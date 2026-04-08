package json

func (c *Client[Db]) ApplyOperation(f func(db *Db) (*Db, error)) error {
	c.m.Lock()
	defer c.m.Unlock()
	db, err := f(&c.db)
	if err != nil {
		return err
	}
	c.db = *db
	c.autosave()
	return nil
}
