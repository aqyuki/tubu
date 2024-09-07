package profile

import (
	"time"
)

type Profile struct {
	// Token is a Discord bot token. It is required.
	Token string

	// Timeout is a duration for Discord API requests.
	Timeout time.Duration
}

func (c Profile) IsValid() bool {
	return c.Token != ""
}
