package controllers

import (
	"fmt"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetBooks func gets all exists books.
// @Description Get all exists books.
// @Summary get all exists books
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /v1/books [get]
func GetCart(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	now := time.Now().Unix()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
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

	// Get all books.
	carts, err := db.GetCartByUser(claims.UserID)
	fmt.Println(carts)
	if err != nil {
		// Return, if books not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"count": 0,
			"books": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(carts),
		"books": carts,
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
func CreateCart(c *fiber.Ctx) error {
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
	cart := &models.Cart{}
	cart.UserID = claims.UserID

	// Check, if received JSON data is valid.
	if err := c.BodyParser(cart); err != nil {
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
	foundedProduct, err := db.GetProductById(cart.Product_id)
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product not found",
		})
	}
	if foundedProduct.Id == 0 {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product not found" + err.Error(),
		})
	}
	if foundedProduct.Quantity < cart.Quantity {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product quantity is not enough",
		})
	}

	foundedCart, err := db.GetCartByUserAndProduct(claims.UserID, cart.Product_id)
	if err != nil {
		// Return status 500 and database connection error.
		if err.Error() != "sql: no rows in result set" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	// Create a new validator for a Book model.
	validate := utils.NewValidator()

	cart.UserID = claims.UserID

	// Validate book fields.
	if err := validate.Struct(cart); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}
	if foundedCart.ID != 0 {
		foundedCart.Quantity += cart.Quantity
		if foundedCart.Quantity > foundedProduct.Quantity {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "product quantity is not enough",
			})
		}
		if err := db.UpdateCart(foundedCart.ID, foundedCart.Quantity); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	} else {
		// Create book by given model.
		if err := db.CreateCart(cart); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}
	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"book":  cart,
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
func UpdateCart(c *fiber.Ctx) error {
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
	cart := &models.Cart{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(cart); err != nil {
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
	foundedCart, err := db.GetCartById(cart.ID)
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
	if foundedCart.UserID == userID {
		// Set initialized default data for book:

		// Create a new validator for a Book model.
		validate := utils.NewValidator()

		// Validate book fields.
		if err := validate.Struct(cart); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}
		foundedProduct, err := db.GetProductById(cart.Product_id)
		if err != nil {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "product not found",
			})
		}
		if foundedProduct.Id == 0 {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "product not found" + err.Error(),
			})
		}
		if foundedProduct.Quantity < cart.Quantity+foundedCart.Quantity {
			// Return status 500 and database connection error.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "product quantity is not enough",
			})
		}

		// Update book by given ID.
		if err := db.UpdateCart(foundedCart.ID, cart.Quantity); err != nil {
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
			"msg":   "permission denied, only the creator can delete his book",
		})
	}
}

// DeleteBook func for deletes book by given ID.
// @Description Delete book by given ID.
// @Summary delete book by given ID
// @Tags Book
// @Accept json
// @Produce json
// @Param id body string true "Book ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/book [delete]
func DeleteCart(c *fiber.Ctx) error {
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
	cart := &models.Cart{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(cart); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Book model.
	validate := utils.NewValidator()

	// Validate book fields.
	if err := validate.StructPartial(cart, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
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
	foundedCart, err := db.GetCartById(cart.ID)
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
	if foundedCart.UserID == userID {
		// Delete book by given ID.
		if err := db.DeleteCartById(foundedCart.ID); err != nil {
			// Return status 500 and error message.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Return status 204 no content.
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, only the creator can delete his book",
		})
	}
}
