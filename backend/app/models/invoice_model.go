package models

type Invoice struct {
	ID           int           `db:"id" json:"id"`
	CreatedAt    string        `db:"created_at" json:"created_at"`
	UpdatedAt    string        `db:"updated_at" json:"updated_at"`
	UserID       int           `db:"user_id" json:"user_id" validate:"required"`
	Email        string        `db:"email" json:"email" `
	Total        int           `db:"total" json:"total" validate:"required"`
	Status       string        `db:"status" json:"status" validate:"required"`
	InvoiceItems []InvoiceItem `db:"invoice_items" json:"invoice_items"`
}

type InvoiceItem struct {
	ID        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	InvoiceID int    `db:"invoice_id" json:"invoice_id" `
	ProductID int    `db:"product_id" json:"product_id" validate:"required"`
	Product   string `json:"product" validate:"required"`
	Quantity  int    `db:"quantity" json:"quantity" validate:"required"`
	Price     int    `db:"price" json:"price" validate:"required"`
}

type PayloadInvoice struct {
	Cart_id int `json:"cart_id"`
}
