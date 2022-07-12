package product

import "time"

type Product struct {
	Id        int
	Name      string
	Desc      string
	SKU       string
	Category  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
