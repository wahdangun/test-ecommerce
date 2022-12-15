package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID         int       `db:"id" json:"id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid"`
	Quantity   int       `db:"quantity" json:"quantity" validate:"required"`
	Product_id int       `db:"product_id" json:"product_id" validate:"required"`
}
