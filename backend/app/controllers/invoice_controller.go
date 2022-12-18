package controllers

import (
	"strconv"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetBook func gets book by given ID or 404 error.
// @Description Get book by given ID.
// @Summary get book by given ID
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.Book
// @Router /v1/book/{id} [get]
func GetInvoiceById(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "order not found",
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	id_int, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Order is not found",
			"book":  nil,
		})
	}
	// Get book by ID.
	book, err := db.GetInvoiceById(id_int)
	if err != nil {
		// Return, if book not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Order is not found",
			"book":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"book":  book,
	})
}

// get invoices by id user
func GetInvoicesByUserId(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	id := claims.UserID
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Order is not found" + err.Error(),
			"book":  nil,
		})
	}
	// Get book by ID.
	book, err := db.GetInvoiceByUser(id)
	if err != nil {
		// Return, if book not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Order is not found",
			"book":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"book":  book,
	})
}

// CreateBook func for creates a new book.
// @Description Create a new book.
// @Summary create a new book
// @Tags Book
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param book_attrs body models.BookAttrs true "Book attributes"
// @Success 200 {object} models.Book
// @Security ApiKeyAuth
// @Router /v1/book [post]
func CreateInvoice(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new Book struct
	payload := []models.PayloadInvoice{}
	invoice := &models.Invoice{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(&payload); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	for _, v := range payload {
		foundedCart, err := db.GetCartById(v.Cart_id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "Cart is not found",
				"book":  nil,
			})
		}
		invoice.InvoiceItems = append(invoice.InvoiceItems, models.InvoiceItem{ProductID: foundedCart.Product_id, Price: foundedCart.Price, Quantity: foundedCart.Quantity, Product: foundedCart.Title})
		invoice.Total += foundedCart.Price * foundedCart.Quantity

	}

	// Create a new validator for a Book model.
	validate := utils.NewValidator()

	// Set initialized default data for book:

	invoice.UserID = claims.UserID
	invoice.Status = "unpaid" // Default status is unpaid.

	// Validate book fields.
	if err := validate.Struct(invoice); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create book by given model.
	if err := db.CreateInvoice(invoice); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	for _, v := range payload {
		db.DeleteCartById(v.Cart_id)
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"book":  invoice,
	})
}

// UpdateBook func for updates book by given ID.
// @Description Update book.
// @Summary update book
// @Tags Book
// @Accept json
// @Produce json
// @Param id body string true "Book ID"
// @Param title body string true "Title"
// @Param author body string true "Author"
// @Param user_id body string true "User ID"
// @Param book_status body integer true "Book status"
// @Param book_attrs body models.BookAttrs true "Book attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/book [put]
func UpdateInvoice(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `book:update` from JWT data of current book.
	credential := claims.Credentials[repository.BookUpdateCredential]

	// Only book creator with `book:update` credential can update his book.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new Book struct
	invoice := &models.Invoice{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(invoice); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if book with given ID is exists.
	foundedInvoice, err := db.GetInvoiceById(invoice.ID)
	if err != nil {
		// Return status 404 and book not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "book with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his book.
	if foundedInvoice.UserID == userID {
		// Set initialized default data for book:

		// Create a new validator for a Book model.
		validate := utils.NewValidator()

		// Validate book fields.
		if err := validate.Struct(invoice); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Update book by given ID.
		if err := db.UpdateInvoice(foundedInvoice.ID, invoice.Status); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 201.
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": false,
			"msg":   nil,
		})
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, invalid user",
		})
	}
}
