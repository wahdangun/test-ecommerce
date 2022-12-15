package queries

import (
	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/jmoiron/sqlx"
)

// BookQueries struct for queries from Book model.
type InvoiceQueries struct {
	*sqlx.DB
}

// GetBooks method for getting all books.
func (q *ProductQueries) GetInvoices() ([]models.Invoice, error) {
	// Define books variable.
	invoices := []models.Invoice{}

	// Define query string.
	query := `SELECT * FROM invoices`

	// Send query to database.
	err := q.Get(&invoices, query)
	if err != nil {
		// Return empty object and error.
		return invoices, err
	}

	// Return query result.
	return invoices, nil
}

// GetBooksByAuthor method for getting all books by given author.
func (q *ProductQueries) GetInvoiceByUser(user string) ([]models.Invoice, error) {
	// Define books variable.
	invoices := []models.Invoice{}

	// Define query string.
	query := `SELECT * FROM invoices left join users on user.id = invoices.user_id WHERE user_id = $1`
	queryItems := `SELECT * FROM invoice_items WHERE invoice_id = $1`

	// Send query to database.
	err := q.Get(&invoices, query, user)
	if err != nil {
		// Return empty object and error.
		return invoices, err
	}
	for i, invoice := range invoices {
		err := q.Get(&invoices[i].InvoiceItems, queryItems, invoice.ID)
		if err != nil {
			// Return empty object and error.
			return invoices, err
		}
	}

	// Return query result.
	return invoices, nil
}

// GetBook method for getting one book by given ID.
func (q *InvoiceQueries) GetInvoiceById(id int) (models.Invoice, error) {
	// Define book variable.
	invoices := models.Invoice{}

	// Define query string.
	query := `SELECT * FROM invoices WHERE id = $1`
	queryItems := `SELECT * FROM invoice_items WHERE invoice_id = $1`

	// Send query to database.
	err := q.Get(&invoices, query, id)
	if err != nil {
		// Return empty object and error.
		return invoices, err
	}
	err = q.Get(&invoices.InvoiceItems, queryItems, invoices.ID)
	if err != nil {
		// Return empty object and error.
		return invoices, err
	}

	// Return query result.
	return invoices, nil
}

// CreateBook method for creating book by given Book object.
func (q *ProductQueries) CreateInvoice(b *models.Invoice) error {
	// Define query string.
	query := `INSERT INTO Invoice VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	updateStock := `UPDATE products SET quantity = quantity - $1 WHERE id = $2`

	// Send query to database.
	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Status, b.Total)
	if err != nil {
		// Return only error.
		return err
	}
	for _, v := range b.InvoiceItems {
		_, err = q.Exec(updateStock, v.Quantity, v.ProductID)
		if err != nil {
			// Return only error.
			return err
		}
		_, err = q.Exec(`INSERT INTO invoice_items VALUES ($1, $2, $3, $4, $5)`, v.ID, v.CreatedAt, v.UpdatedAt, v.InvoiceID, v.ProductID, v.Quantity, v.Price)
		if err != nil {
			// Return only error.
			return err
		}
	}
	// This query returns nothing.
	return nil
}

// UpdateBook method for updating book by given Book object.
func (q *ProductQueries) UpdateInvoice(id int, status string) error {
	// Define query string.
	query := `UPDATE invoices SET updated_at = $2, status = $3 WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id, status)
	if err != nil {
		// Return only error.
		return err
	}
	if status == "cancel" {
		queryItems := `SELECT * FROM invoice_items WHERE invoice_id = $1`
		var items []models.InvoiceItem
		err = q.Get(&items, queryItems, id)
		if err != nil {
			// Return empty object and error.
			return err
		}
		updateStock := `UPDATE products SET quantity = quantity + $1 WHERE id = $2`
		for _, v := range items {
			_, err = q.Exec(updateStock, v.Quantity, v.ProductID)
			if err != nil {
				// Return only error.
				return err
			}
		}
	}

	// This query returns nothing.
	return nil
}
