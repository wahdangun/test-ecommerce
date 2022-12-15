package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	ID           int           `db:"id" json:"id"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at" json:"updated_at"`
	UserID       uuid.UUID     `db:"user_id" json:"user_id" validate:"required,uuid"`
	Total        int           `db:"total" json:"total" validate:"required"`
	Status       string        `db:"status" json:"status" validate:"required"`
	InvoiceItems []InvoiceItem `db:"invoice_items" json:"invoice_items" validate:"required,dive"`
}

type InvoiceItem struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	InvoiceID int       `db:"invoice_id" json:"invoice_id" validate:"required"`
	ProductID int       `db:"product_id" json:"product_id" validate:"required"`
	Product   string    `json:"product" validate:"required"`
	Quantity  int       `db:"quantity" json:"quantity" validate:"required"`
	Price     int       `db:"price" json:"price" validate:"required"`
}

// Value make the ProductAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b InvoiceItem) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the ProductAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *InvoiceItem) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}
