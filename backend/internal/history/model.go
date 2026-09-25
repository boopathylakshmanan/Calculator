package history

import "time"

type History struct {
	ID         int64     `json:"id"`
	Expression string    `json:"expression"`
	Result     string    `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}
