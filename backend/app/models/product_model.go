package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Product struct to describe product object.
type Product struct {
	ID            uuid.UUID    `db:"id" json:"id" validate:"required,uuid"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at" json:"updated_at"`
	UserID        uuid.UUID    `db:"user_id" json:"user_id" validate:"required,uuid"`
	Title         string       `db:"title" json:"title" validate:"required,lte=255"`
	Price         int          `db:"price" json:"price" validate:"required"`
	Quantity      int          `db:"quantity" json:"quantity" validate:"required"`
	ProductStatus int          `db:"product_status" json:"product_status" validate:"required,len=1"`
	ProductAttrs  ProductAttrs `db:"product_attrs" json:"product_attrs" validate:"required,dive"`
}

// ProductAttrs struct to describe product attributes.
type ProductAttrs struct {
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"min=1,max=10"`
}

// Value make the ProductAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b ProductAttrs) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the ProductAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *ProductAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}