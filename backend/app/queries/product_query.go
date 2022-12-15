package queries

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/jmoiron/sqlx"
)

// BookQueries struct for queries from Book model.
type ProductQueries struct {
	*sqlx.DB
}

// GetBooks method for getting all books.
func (q *ProductQueries) GetProducts() ([]models.Product, error) {
	// Define books variable.
	products := []models.Product{}

	// Define query string.
	query := `SELECT * FROM products`

	// Send query to database.
	err := q.Get(&products, query)
	if err != nil {
		// Return empty object and error.
		return products, err
	}

	// Return query result.
	return products, nil
}

// GetBooksByAuthor method for getting all books by given author.
func (q *ProductQueries) GetProductByTitle(title string) ([]models.Product, error) {
	// Define books variable.
	products := []models.Product{}

	// Define query string.
	query := `SELECT * FROM products WHERE title = $1`

	// Send query to database.
	err := q.Get(&products, query, title)
	if err != nil {
		// Return empty object and error.
		return products, err
	}

	// Return query result.
	return products, nil
}

// GetBook method for getting one book by given ID.
func (q *ProductQueries) GetProductById(id int) (models.Product, error) {
	// Define book variable.
	products := models.Product{}

	// Define query string.
	query := `SELECT * FROM products WHERE id = $1`

	// Send query to database.
	err := q.Get(&products, query, id)
	if err != nil {
		// Return empty object and error.
		return products, err
	}

	// Return query result.
	return products, nil
}

// CreateBook method for creating book by given Book object.
func (q *ProductQueries) CreateProduct(b *models.Product) error {
	// Define query string.
	query := `INSERT INTO products VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Send query to database.
	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Price, b.ProductStatus, b.ProductAttrs, b.Quantity)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateBook method for updating book by given Book object.
func (q *ProductQueries) UpdateProduct(id int, b *models.Product) error {
	// Define query string.
	query := `UPDATE books SET updated_at = $2, title = $3, price = $4, quantity = $5, product_status = $6, product_attrs = $7 WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id, b.UpdatedAt, b.Title, b.Price, b.Quantity, b.ProductStatus, b.ProductAttrs)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBook method for delete book by given ID.
func (q *BookQueries) DeleteProduct(id int) error {
	// Define query string.
	query := `DELETE FROM products WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
