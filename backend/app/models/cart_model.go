package models

type Cart struct {
	ID          int    `db:"id" json:"id"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
	UserID      int    `db:"user_id" json:"user_id" validate:"required"`
	Title       string `db:"title" json:"title"`
	Quantity    int    `db:"quantity" json:"quantity" validate:"required"`
	Price       int    `db:"price" json:"price"`
	Product_id  int    `db:"product_id" json:"product_id" validate:"required"`
	Cart_status int    `db:"cart_status" json:"cart_status"`
}
