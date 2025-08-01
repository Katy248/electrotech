package orders

import (
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/repository/users"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	Products []OrderProductRequest `json:"products" binding:"required,min=1"`
}

type OrderProductRequest struct {
	ProductId string `json:"id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

func checkUserHasCompanyData(user users.User) bool {
	return user.CompanyName.Valid &&
		user.CompanyAddress.Valid &&
		user.PositionInCompany.Valid && user.CompanyInn.Valid
}

func CreateOrderHandler(orderRepo *orders.Queries, userRepo *users.Queries, catalogRepo *catalog.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста (предполагаем, что middleware аутентификации уже добавил его)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, существует ли пользователь
		user, err := userRepo.GetById(c.Request.Context(), userID.(int64))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if !checkUserHasCompanyData(user) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user has no company data"})
			return
		}

		// Конвертируем продукты запроса в структуры репозитория
		var products []orders.AddOrderProductParams
		for _, p := range req.Products {
			price, err := catalogRepo.GetProductPrice(p.ProductId)
			if err == catalog.ErrNotFound {
				log.Printf("Error getting product price: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return

			} else if err != nil {
				log.Printf("Something bad with product product: %v", err)
				continue

			}

			name, err := catalogRepo.GetProductName(p.ProductId)
			if err != nil {
				log.Printf("something bad with product name: %v", err)
				continue
			}

			products = append(products, orders.AddOrderProductParams{
				ProductName: name,
				Quantity:    int64(p.Quantity),
				Price:       float64(price),
			})
		}

		if len(products) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no products in order"})
			return
		}

		// Создаем заказ
		order, err := orderRepo.InsertOrder(c.Request.Context(), orders.InsertOrderParams{
			UserID:       userID.(int64),
			CreationDate: time.Now(),
		})
		if err != nil {
			log.Printf("Error creating order: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
			return
		}

		for _, p := range products {
			err := orderRepo.AddOrderProduct(c.Request.Context(), orders.AddOrderProductParams{
				OrderID:     order.ID,
				ProductName: p.ProductName,
				Quantity:    p.Quantity,
				Price:       p.Price,
			})
			if err != nil {
				log.Printf("Error adding product to order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product to order"})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "order created successfully",
			"orderId": order.ID,
		})
	}
}

func GetUserOrdersHandler(orderRepo *orders.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		orders, err := orderRepo.GetOrders(c.Request.Context(), userID.(int64))
		if err != nil {
			log.Printf("Error getting user orders: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
			return
		}

		responseOrders := []ResponseOrder{}

		for _, order := range orders {
			products, err := orderRepo.GetOrderProducts(c.Request.Context(), order.ID)
			if err != nil {
				log.Printf("Error getting order: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
				return
			}
			responseOrder := newResponseOrder(order, products)
			responseOrders = append(responseOrders, responseOrder)
		}

		// Получаем заказы пользователя
		c.JSON(http.StatusOK, gin.H{"orders": responseOrders})
	}
}

func newResponseOrder(order orders.Order, products []orders.OrderProduct) ResponseOrder {
	response := ResponseOrder{
		ID:        order.ID,
		CreatedAt: order.CreationDate.String(),
	}

	for _, p := range products {
		response.Products = append(response.Products, ResponseOrderProduct{
			Name:     p.ProductName,
			Id:       p.ID,
			Quantity: p.Quantity,
			Price:    p.ProductPrice,
		})
	}
	return response
}

type ResponseOrder struct {
	ID        int64                  `json:"id"`
	CreatedAt string                 `json:"createdAt"`
	Products  []ResponseOrderProduct `json:"products"`
}

type ResponseOrderProduct struct {
	Name     string  `json:"productName"`
	Id       int64   `json:"productId"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}
