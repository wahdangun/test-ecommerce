package controllers

import (
	"strconv"
	"time"

	models "github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
)

// Getinvoice func gets invoice by given ID or 404 error.
// @Description Get invoice by given ID.
// @Summary get invoice by given ID
// @Tags invoice
// @Accept json
// @Produce json
// @Param id path string true "invoice ID"
// @Success 200 {object} models.Invoice
// @Security ApiKeyAuth
// @Router /v1/invoice/{id} [get]
func GetInvoiceById(c *fiber.Ctx) error {
	// Catch invoice ID from URL.
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
			"error":   true,
			"msg":     "Order is not found",
			"invoice": nil,
		})
	}
	// Get invoice by ID.
	invoice, err := db.GetInvoiceById(id_int)
	if err != nil {
		// Return, if invoice not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "Order is not found",
			"invoice": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"invoice": invoice,
	})
}

// GetInvoicesByUserId func gets all exists invoice by user id in jwt claims.
// @Description Get all exists invoice.
// @Summary get all exists invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Success 200 {array} models.Invoice
// @Security ApiKeyAuth
// @Router /v1/invoice [get]
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

	// Set expiration time from JWT data of current invoice.
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
			"error":   true,
			"msg":     "Order is not found" + err.Error(),
			"invoice": nil,
		})
	}
	// Get invoice by ID.
	invoice, err := db.GetInvoiceByUser(id)
	if err != nil {
		// Return, if invoice not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "Order is not found",
			"invoice": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"invoice": invoice,
	})
}

// Createinvoice func for creates a new invoice.
// @Description Create a new invoice.
// @Summary create a new invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Param cart_id body models.PayloadInvoice true "Cart IDs"
// @Success 200 {object} models.Invoice
// @Security ApiKeyAuth
// @Router /v1/invoice [post]
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

	// Set expiration time from JWT data of current invoice.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new invoice struct
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
				"error":   true,
				"msg":     "Cart is not found",
				"invoice": nil,
			})
		}
		invoice.InvoiceItems = append(invoice.InvoiceItems, models.InvoiceItem{ProductID: foundedCart.Product_id, Price: foundedCart.Price, Quantity: foundedCart.Quantity, Product: foundedCart.Title})
		invoice.Total += foundedCart.Price * foundedCart.Quantity

	}

	// Create a new validator for a invoice model.
	validate := utils.NewValidator()

	// Set initialized default data for invoice:

	invoice.UserID = claims.UserID
	invoice.Status = "unpaid" // Default status is unpaid.

	// Validate invoice fields.
	if err := validate.Struct(invoice); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create invoice by given model.
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
		"error":   false,
		"msg":     nil,
		"invoice": invoice,
	})
}

// Updateinvoice func for updates invoice by given ID.
// @Description Update invoice.
// @Summary update invoice
// @Tags invoice
// @Accept json
// @Produce json
// @Param id body string true "invoice ID"
// @Param title body string true "Title"
// @Param user_id body string true "User ID"
// @Param invoice_status body integer true "invoice status"
// @Param invoice_attrs body models.InvoiceItem true "invoice items"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/invoice [put]
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

	// Set expiration time from JWT data of current invoice.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `invoice:update` from JWT data of current invoice.
	credential := claims.Credentials[repository.InvoiceUpdateCredential]

	// Only invoice creator with `invoice:update` credential can update his invoice.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new invoice struct
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

	// Checking, if invoice with given ID is exists.
	foundedInvoice, err := db.GetInvoiceById(invoice.ID)
	if err != nil {
		// Return status 404 and invoice not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "invoice with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his invoice.
	if foundedInvoice.UserID == userID {
		// Set initialized default data for invoice:

		// Create a new validator for a invoice model.
		validate := utils.NewValidator()

		// Validate invoice fields.
		if err := validate.Struct(invoice); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Update invoice by given ID.
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
