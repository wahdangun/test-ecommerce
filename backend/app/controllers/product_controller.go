package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/create-go-app/fiber-go-template/pkg/repository"
	"github.com/create-go-app/fiber-go-template/pkg/utils"
	"github.com/create-go-app/fiber-go-template/platform/database"
	"github.com/gofiber/fiber/v2"
)

// Getproducts func gets all exists products.
// @Description Get all exists products.
// @Summary get all exists products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Router /v1/products [get]
func GetProducts(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all products.
	products, err := db.GetProducts()
	if err != nil {
		// Return, if products not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "products were not found" + err.Error(),
			"count":    0,
			"products": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"count":   len(products),
		"product": products,
	})
}

// Getproduct func gets product by given ID or 404 error.
// @Description Get product by given ID.
// @Summary get product by given ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "product ID"
// @Success 200 {object} models.Product
// @Router /v1/product/{id} [get]
func GetProduct(c *fiber.Ctx) error {
	// Catch product ID from URL.
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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

	// Get product by ID.
	product, err := db.GetProductById(id)
	if err != nil {
		// Return, if product not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "product with the given ID is not found",
			"product": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"product": product,
	})
}

// Createproduct func for creates a new product.
// @Description Create a new product.
// @Summary create a new product
// @Tags product
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param price body int true "Price"
// @Param quantity body int true "Quantity"
// @Param product_attrs body models.ProductAttrs true "product attributes"
// @Success 200 {object} models.Product
// @Security ApiKeyAuth
// @Router /v1/product [post]
func CreateProduct(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current product.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	//Set credential `product:create` from JWT data of current product.
	credential := claims.Credentials[repository.ProductCreateCredential]

	//	Only user with `product:create` credential can create a new product.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new product struct
	product := &models.Product{}
	fmt.Println("product", product)
	// Check, if received JSON data is valid.
	if err := c.BodyParser(product); err != nil {
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

	// Create a new validator for a product model.
	validate := utils.NewValidator()

	product.User_id = claims.UserID
	product.ProductStatus = 1 // 0 == draft, 1 == active

	// Validate product fields.
	if err := validate.Struct(product); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create product by given model.
	if err := db.CreateProduct(product); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"product": product,
	})
}

// Updateproduct func for updates product by given ID.
// @Description Update product.
// @Summary update product
// @Tags product
// @Accept json
// @Produce json
// @Param title body string true "Title"
// @Param price body int true "Price"
// @Param quantity body int true "Quantity"
// @Param product_status body integer true "product status"
// @Param product_attrs body models.ProductAttrs true "product attributes"
// @Success 202 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/product [put]
func UpdateProduct(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current product.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Set credential `product:update` from JWT data of current product.
	credential := claims.Credentials[repository.ProductUpdateCredential]

	// Only product creator with `product:update` credential can update his product.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
		})
	}

	// Create new product struct
	product := &models.Product{}
	product.User_id = claims.UserID
	// Check, if received JSON data is valid.
	if err := c.BodyParser(product); err != nil {
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

	// Checking, if product with given ID is exists.
	foundedProduct, err := db.GetProductById(product.Id)
	if err != nil {
		// Return status 404 and product not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his product.
	if foundedProduct.User_id == userID {
		// Set initialized default data for product:

		// Create a new validator for a product model.
		validate := utils.NewValidator()

		// Validate product fields.
		if err := validate.Struct(product); err != nil {
			// Return, if some fields are not valid.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   utils.ValidatorErrors(err),
			})
		}

		// Update product by given ID.
		if err := db.UpdateProduct(foundedProduct.Id, product); err != nil {
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
			"msg":   "permission denied, only the creator can delete his product",
		})
	}
}

// Deleteproduct func for deletes product by given ID.
// @Description Delete product by given ID.
// @Summary delete product by given ID
// @Tags product
// @Accept json
// @Produce json
// @Param id body string true "product ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/product [delete]
func DeleteProduct(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current product.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}
	//Set credential `product:delete` from JWT data of current product.
	credential := claims.Credentials[repository.ProductDeleteCredential]

	//	Only user with `product:delete` credential can delete a  product.
	if !credential {
		// Return status 403 and permission denied error message.
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": true,
			"msg":   "permission denied, check credentials of your token",
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

	// Checking, if product with given ID is exists.
	foundedProduct, err := db.GetProductById(id)
	if err != nil {
		// Return status 404 and product not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "product with this ID not found",
		})
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his product.
	if foundedProduct.User_id == userID {
		// Delete product by given ID.
		if err := db.DeleteProduct(foundedProduct.Id); err != nil {
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
			"msg":   "permission denied, only the creator can delete his product",
		})
	}
}
