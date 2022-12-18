package queries

import (
	"fmt"

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
	query := `SELECT id, created_at, updated_at, title, quantity, price, product_status, product_attrs, user_id FROM products`
	// Send query to database.
	err := q.Select(&products, query)
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
	query := `SELECT * FROM products WHERE title = ?`

	// Send query to database.
	err := q.Select(&products, query, title)
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
	product := models.Product{}

	// Define query string.
	query := `SELECT * FROM products WHERE id = ?`

	// Send query to database.
	err := q.Get(&product, query, id)
	if err != nil {
		// Return empty object and error.
		fmt.Println(err.Error(), query)
		return product, err
	}

	// Return query result.
	return product, nil
}

// CreateBook method for creating book by given Book object.
func (q *ProductQueries) CreateProduct(b *models.Product) error {
	// Define query string.
	fmt.Println(b)
	query := `INSERT INTO products (user_id, title, price, product_status, quantity, product_attrs) VALUES (?,?,?,?,?,?)`

	// Send query to database.
	_, err := q.Exec(query, b.User_id, b.Title, b.Price, b.ProductStatus, b.Quantity, b.ProductAttrs)
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
	query := `UPDATE products SET  title = ?, price = ?, quantity = ?, product_status = ?, product_attrs = ? WHERE id = ?`

	// Send query to database.
	_, err := q.Exec(query, b.Title, b.Price, b.Quantity, b.ProductStatus, b.ProductAttrs, id)
	if err != nil {
		// Return only error.
		fmt.Println(err.Error(), query)
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBook method for delete book by given ID.
func (q *ProductQueries) DeleteProduct(id int) error {
	// Define query string.
	query := `DELETE FROM products WHERE id = ?`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
