package json

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Client[Db any] struct {
	Path          string
	emptyDbObject []byte
	db            Db
	m             sync.Mutex
}

func New[Db any](path string) (*Client[Db], error) {
	var empty Db
	emptyDbObject, err := json.Marshal(empty)
	if err != nil {
		return nil, fmt.Errorf("marshalling empty database object failed: %w", err)
	}

	c := &Client[Db]{
		Path:          path,
		m:             sync.Mutex{},
		emptyDbObject: emptyDbObject,
	}
	if !c.dbExists() {
		err := c.writeDb()
		if err != nil {
			return nil, fmt.Errorf("failed to write default db: %w", err)
		}
	}
	if err := c.readdb(); err != nil {
		return nil, fmt.Errorf("failed to load the database: %w", err)
	}
	return c, nil
}

func (c *Client[Db]) dbExists() bool {
	_, err := os.Stat(c.Path)
	return os.IsExist(err) || err == nil
}

func (c *Client[Db]) writeDb() error {
	err := os.WriteFile(c.Path, []byte(`{"users":[]}`), 0o644)
	if err != nil {
		return fmt.Errorf("failed to create default db: %w", err)
	}
	return nil
}
