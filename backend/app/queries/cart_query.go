package queries

import (
	"fmt"

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
	query := `SELECT carts.id, carts.user_id,if(carts.quantity > products.quantity,products.quantity, carts.quantity) as quantity, products.title, products.price, carts.product_id, carts.created_at, carts.updated_at FROM carts left join products on products.id = carts.product_id left join users on users.id = carts.user_id WHERE users.id = ?`

	// Send query to database.
	err := q.Select(&carts, query, id)
	if err != nil {
		// Return empty object and error.
		fmt.Println(err.Error())
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
	query := `SELECT carts.id, carts.user_id,if(carts.quantity > products.quantity,products.quantity, carts.quantity) as quantity, products.title, products.price, carts.product_id, carts.created_at, carts.updated_at FROM carts left join products on products.id = carts.product_id left join users on users.id = carts.user_id WHERE carts.id = ?`

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
	fmt.Println(b)
	query := `INSERT INTO carts (user_id, quantity, product_id) VALUES (?,?,?)`

	// Send query to database.
	_, err := q.Exec(query, b.UserID, b.Quantity, b.Product_id)
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
	query := `UPDATE carts SET  quantity = ? WHERE id = ?`

	// Send query to database.
	_, err := q.Exec(query, quantity, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBook method for delete book by given ID.
func (q *BookQueries) DeleteCartById(id int) error {
	// Define query string.
	query := `DELETE FROM carts WHERE id = ?`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// get cart by user id and product id
func (q *CartQueries) GetCartByUserAndProduct(user_id, product_id int) (models.Cart, error) {
	// Define book variable.
	carts := models.Cart{}

	// Define query string.
	query := `SELECT * FROM carts WHERE user_id = ? and product_id = ?`

	// Send query to database.
	err := q.Get(&carts, query, user_id, product_id)
	if err != nil {
		// Return empty object and error.
		return carts, err
	}

	// Return query result.
	return carts, nil
}
