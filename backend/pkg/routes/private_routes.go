package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)           // create a new book
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

	// Routes for PUT method:
	route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// Routes for DELETE method:
	route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID

	route.Post("/product", middleware.JWTProtected(), controllers.CreateProduct)       // create a new product
	route.Put("/product", middleware.JWTProtected(), controllers.UpdateProduct)        // update product by ID
	route.Delete("/product/:id", middleware.JWTProtected(), controllers.DeleteProduct) // delete product by ID

	route.Get("/cart", controllers.GetCart)           // get list of all carts
	route.Put("/cart", controllers.UpdateCart)        // update cart by ID
	route.Delete("/cart/:id", controllers.DeleteCart) // delete cart by ID
	route.Post("/cart", controllers.CreateCart)       // create a new cart

	route.Post("/invoice", controllers.CreateInvoice)      // create a new invoice
	route.Get("/invoice", controllers.GetInvoicesByUserId) // get list of all invoices
	route.Put("/invoice", controllers.UpdateInvoice)       // update invoice by ID
}
