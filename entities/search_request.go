package entities

import (
	"time"
)

type SearchRequest struct {
	OrigPoint      Location
	DestPoint      Location
	CreatedAt      time.Time
	WaitingSeconds uint16
}
