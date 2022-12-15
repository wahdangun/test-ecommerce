package queries

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/jmoiron/sqlx"
)

// BookQueries struct for queries from Book model.
type CartQueries struct {
	*sqlx.DB
}

// GetBooksByAuthor method for getting all books by given author.
func (q *CartQueries) GetCartByUser(id int) ([]models.Cart, error) {
	// Define books variable.
	carts := []models.Cart{}

	// Define query string.
	query := `SELECT * FROM carts left join user on user.id = cart.user_id WHERE user.id = ?`

	// Send query to database.
	err := q.Get(&carts, query, id)
	if err != nil {
		// Return empty object and error.
		return carts, err
	}

	// Return query result.
	return carts, nil
}

// GetCartById method for getting one cart item by given ID.
func (q *CartQueries) GetCartById(id int) (models.Cart, error) {
	// Define book variable.
	carts := models.Cart{}

	// Define query string.
	query := `SELECT * FROM carts WHERE id = $1`

	// Send query to database.
	err := q.Get(&carts, query, id)
	if err != nil {
		// Return empty object and error.
		return carts, err
	}

	// Return query result.
	return carts, nil
}

// CreateBook method for creating book by given Book object.
func (q *CartQueries) CreateCart(b *models.Cart) error {
	// Define query string.
	query := `INSERT INTO carts VALUES ($1, $2, $3, $4, $5, $6, )`

	// Send query to database.
	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Quantity, b.Product_id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateBook method for updating book by given Book object.
func (q *ProductQueries) UpdateCart(id, quantity int) error {
	// Define query string.
	query := `UPDATE carts SET updated_at = $2, quantity = $3 WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id, quantity)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBook method for delete book by given ID.
func (q *BookQueries) DeleteCart(id int) error {
	// Define query string.
	query := `DELETE FROM carts WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
