package contact

import "time"

type Contact struct {
	UserID    string
	ID        string
	Phone     string
	FirstName string
	LastName  string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
