package dto

import "time"

type Slot struct {
	SlotID    string
	Location  string
	DateTime  time.Time
	Available bool
}

type CreateSlotRequest struct {
	SlotID    string
	Location  string
	DateTime  time.Time
	Available bool
}
