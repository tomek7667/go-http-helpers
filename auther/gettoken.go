package auther

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tomek7667/go-http-helpers/utils"
)

func (c *Client[User]) GetToken(user User) (string, error) {
	var umap map[string]any
	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to json.marshal the user: %w", err)
	}

	err = json.Unmarshal(b, &umap)
	if err != nil {
		return "", fmt.Errorf("failed to json.unmarshal the user: %w", err)
	}
	umap["token_created"] = time.Now()

	return utils.JwtEncode(umap, c.Secret)
}
