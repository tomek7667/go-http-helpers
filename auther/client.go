package auther

import (
	"context"
)

type HasID interface {
	GetID() string
}

type Dber[User HasID] interface {
	GetUser(ctx context.Context, id string) (User, error)
}

type Client[User HasID] struct {
	Secret string
	Dber   Dber[User]
}

func New[User HasID](secret string, dber Dber[User]) *Client[User] {
	c := &Client[User]{
		Secret: secret,
		Dber:   dber,
	}

	return c
}
