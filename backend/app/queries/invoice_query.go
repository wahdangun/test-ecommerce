package queries

import (
	"fmt"

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
	query := `SELECT invoices.created_at, invoices.updated_at, invoices.invoice_status, users.email FROM invoices left join users on user.id = invoices.user_id `

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
func (q *ProductQueries) GetInvoiceByUser(user int) ([]models.Invoice, error) {
	// Define books variable.
	invoices := []models.Invoice{}
	invoiceItems := []models.InvoiceItem{}

	// Define query string.
	query := `SELECT invoices.id, invoices.created_at, invoices.updated_at, invoice_status as status, invoices.total, users.email FROM invoices left join users on users.id = invoices.user_id WHERE user_id = ?`
	queryItems := `SELECT invoice_items.id, products.title as product, invoice_items.created_at, invoice_items.updated_at, invoice_items.invoice_id,invoice_items.product_id,  invoice_items.price, invoice_items.quantity FROM invoice_items left join products on products.id = invoice_items.product_id WHERE invoice_items.invoice_id = ?`

	// Send query to database.
	err := q.Select(&invoices, query, user)
	if err != nil {

		// Return empty object and error.
		return invoices, err
	}

	for i, invoice := range invoices {
		fmt.Println(i)
		fmt.Println(invoice.ID)
		err := q.Select(&invoiceItems, queryItems, invoice.ID)
		if err != nil {
			fmt.Println(err.Error())
			// Return empty object and error.
			return invoices, err
		}
		invoices[i].InvoiceItems = invoiceItems
	}
	fmt.Println(invoices)
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
	query := `INSERT INTO invoices  (user_id, invoice_status, total) VALUES (?,?, ?)`
	updateStock := `UPDATE products SET quantity = quantity - ? WHERE id = ?`

	// Send query to database.
	result, err := q.Exec(query, b.UserID, b.Status, b.Total)
	if err != nil {
		// Return only error.
		return err
	}
	id, err := result.LastInsertId()
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
		_, err = q.Exec(`INSERT INTO invoice_items (invoice_id, product_id, quantity, price) VALUES (?,?,?,?)`, id, v.ProductID, v.Quantity, v.Price)
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
	query := `UPDATE invoices SET  status = ? WHERE id = ?`

	// Send query to database.
	_, err := q.Exec(query, status, id)
	if err != nil {
		// Return only error.
		return err
	}
	if status == "canceled" {
		queryItems := `SELECT * FROM invoice_items WHERE invoice_id = ?`
		var items []models.InvoiceItem
		err = q.Get(&items, queryItems, id)
		if err != nil {
			// Return empty object and error.
			return err
		}
		updateStock := `UPDATE products SET quantity = quantity + ? WHERE id = ?`
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

// function to get alll invoices item
func (q *ProductQueries) GetInvoiceItems(id int) ([]models.InvoiceItem, error) {
	// Define books variable.
	invoiceItems := []models.InvoiceItem{}

	// Define query string.
	query := `SELECT * FROM invoice_items where invoice_id = ?`

	// Send query to database.
	err := q.Get(&invoiceItems, query, id)
	if err != nil {
		// Return empty object and error.
		return invoiceItems, err
	}

	// Return query result.
	return invoiceItems, nil
}
