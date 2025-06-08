package tariff

import "time"

type Tariff struct {
	Id          int64
	Name        string
	Description string
	Price       float64
	Currency    string
	DurationMin int64
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
