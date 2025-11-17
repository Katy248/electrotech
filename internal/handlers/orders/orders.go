package orders

import (
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	Products []OrderProductRequest `json:"products" binding:"required,min=1"`
}

type OrderProductRequest struct {
	ProductId string `json:"id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

func CreateOrderHandler(catalogRepo *catalog.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста (предполагаем, что middleware аутентификации уже добавил его)
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed to bind request")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := users.ByID(userID.(int64))
		if err != nil {
			log.Error("User not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		products := []models.OrderProduct{}
		for _, p := range req.Products {

			product, err := catalogRepo.GetProduct(p.ProductId)
			if err != nil {
				log.Error("Failed getting product price, product not found", "productId", p.ProductId, "error", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return
			}

			products = append(products,
				models.OrderProduct{
					ProductName:  product.Name,
					Quantity:     int64(p.Quantity),
					ProductPrice: float64(product.Price),
					ProductID:    p.ProductId,
				},
			)

		}
		order, err := orders.New(user, products)
		if err != nil {
			log.Error("Failed creating order", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
			return
		}

		// if err := order.SetUser(user); err != nil {
		// 	log.Error("Can't set user as order creator", "user", user, "error", err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "user can't create orders"})
		// 	return
		// }

		go sendEmail(*order)

		c.JSON(http.StatusCreated, gin.H{
			"message": "order created successfully",
			"orderId": order.ID,
		})
	}
}

func GetUserOrdersHandler(catalogRepo *catalog.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		orders, err := orders.GetOrders(userID.(int64))
		if err != nil {
			log.Error("Failed getting user orders", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
			return
		}

		for _, o := range orders {
			for _, p := range o.OrderProducts {
				product, err := catalogRepo.GetProduct(p.ProductID)
				if err != nil {
					log.Error("Failed getting product", "error", err, "productID", p.ProductID)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
					continue
				}
				p.ImagePath = product.ImagePath
			}

		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}
