package controllers

import (
	"fmt"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
)

// Getcart func gets all exists cart.
// @Description Get all exists cart.
// @Summary get all exists cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {array} models.Cart
// @Router /v1/cart [get]
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

	// Set expiration time from JWT data of current cart.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Get all cart.
	carts, err := db.GetCartByUser(claims.UserID)
	fmt.Println(carts)
	if err != nil {
		// Return, if cart not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"count": 0,
			"cart":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(carts),
		"cart":  carts,
	})
}

// Createcart func for creates a new cart.
// @Description Create a new cart.
// @Summary create a new cart
// @Tags cart
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param user_id body string true "User ID"
// @Success 200 {object} models.Cart
// @Security ApiKeyAuth
// @Router /v1/cart [post]
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

	// Set expiration time from JWT data of current cart.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new cart struct
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

	// Create a new validator for a cart model.
	validate := utils.NewValidator()

	cart.UserID = claims.UserID

	// Validate cart fields.
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
		// Create cart by given model.
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
		"cart":  cart,
	})
}

// Updatecart func for updates cart by given ID.
// @Description Update cart.
// @Summary update cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id body string true "cart ID"
// @Param title body string true "Title"
// @Param user_id body string true "User ID"
// @Param cart_status body integer true "cart status"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/cart [put]
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

	// Set expiration time from JWT data of current cart.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new cart struct
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

	// Checking, if cart with given ID is exists.
	foundedCart, err := db.GetCartById(cart.ID)
	if err != nil {
		// Return status 404 and cart not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "cart with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his cart.
	if foundedCart.UserID == userID {
		// Set initialized default data for cart:

		// Create a new validator for a cart model.
		validate := utils.NewValidator()

		// Validate cart fields.
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

		// Update cart by given ID.
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
			"msg":   "permission denied, only the creator can delete his cart",
		})
	}
}

// Deletecart func for deletes cart by given ID.
// @Description Delete cart by given ID.
// @Summary delete cart by given ID
// @Tags cart
// @Accept json
// @Produce json
// @Param id body string true "cart ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/cart [delete]
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

	// Set expiration time from JWT data of current cart.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new cart struct
	cart := &models.Cart{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(cart); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a cart model.
	validate := utils.NewValidator()

	// Validate cart fields.
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

	// Checking, if cart with given ID is exists.
	foundedCart, err := db.GetCartById(cart.ID)
	if err != nil {
		// Return status 404 and cart not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "cart with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his cart.
	if foundedCart.UserID == userID {
		// Delete cart by given ID.
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
			"msg":   "permission denied, only the creator can delete his cart",
		})
	}
}
