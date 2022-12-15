package models

type Cart struct {
	ID         int    `db:"id" json:"id"`
	CreatedAt  string `db:"created_at" json:"created_at"`
	UpdatedAt  string `db:"updated_at" json:"updated_at"`
	UserID     int    `db:"user_id" json:"user_id" validate:"required"`
	Quantity   int    `db:"quantity" json:"quantity" validate:"required"`
	Product_id int    `db:"product_id" json:"product_id" validate:"required"`
}
