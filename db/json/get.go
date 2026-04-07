package json

func (c *Client[Db]) Get() *Db {
	return &c.db
}
