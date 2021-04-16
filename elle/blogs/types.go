package blogs

import "time"

type Info struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	ContentID   int64
	Price int64
}

type Content struct {
	ID   int64
	Text string
}
