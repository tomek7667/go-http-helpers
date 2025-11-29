package auther

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tomek7667/go-http-helpers/utils"
)

func (c *Client[User]) GetUserFromToken(ctx context.Context, token string) (*User, error) {
	data, err := utils.JwtVerify(token, c.Secret)
	if err != nil {
		return nil, fmt.Errorf("jwt verify failed in getuserfromtoken failed: %w", err)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal in getuserfromtoken failed: %w", err)
	}

	var user User
	err = json.Unmarshal(dataBytes, &user)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal in getuserfromtoken failed: %w", err)
	}

	user, err = c.Dber.GetUser(ctx, user.GetID())
	if err != nil {
		return nil, fmt.Errorf("retrieving the user from the database failed: %w", err)
	}

	return &user, nil
}
