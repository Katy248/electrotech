package orders

import (
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"electrotech/internal/users"
	"fmt"
	"net/http"
	"time"

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

func CreateOrderHandler(orderRepo *orders.Queries, catalogRepo *catalog.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста (предполагаем, что middleware аутентификации уже добавил его)
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		orderModel := Order{}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("Failed to bind request")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, существует ли пользователь
		user, err := users.ByID(userID.(int64))
		if err != nil {
			log.Error("User not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if !user.CompanyData().DataFilled() {
			log.Error("User has no required company data")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user has no company data"})
			return
		}

		// Конвертируем продукты запроса в структуры репозитория
		var products []orders.AddOrderProductParams
		for _, p := range req.Products {
			price, err := catalogRepo.GetProductPrice(p.ProductId)
			if err == catalog.ErrNotFound {
				log.Error("Failed getting product price", "error", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
				return

			} else if err != nil {
				log.Warn("Something bad with order item, skipping it", "item", p, "error", err)
				continue

			}

			name, err := catalogRepo.GetProductName(p.ProductId)
			if err != nil {
				log.Warn("Can't get product name", "productId", p.ProductId, "error", err)
				continue
			}

			products = append(products, orders.AddOrderProductParams{
				ProductName: name,
				Quantity:    int64(p.Quantity),
				Price:       float64(price),
				ProductID:   p.ProductId,
			})

			orderModel.Products = append(orderModel.Products, Product{
				ID:       p.ProductId,
				Quantity: int64(p.Quantity),
				Price:    float64(price),
				Name:     name,
			})
		}

		if len(products) == 0 {
			log.Error("No products in order, counts only valid items", "products", products, "itemsInRequest", len(req.Products))
			c.JSON(http.StatusBadRequest, gin.H{"error": "no products in order"})
			return
		}

		// Создаем заказ
		order, err := orderRepo.InsertOrder(c.Request.Context(), orders.InsertOrderParams{
			UserID:       userID.(int64),
			CreationDate: time.Now(),
		})
		if err != nil {
			log.Error("Failed creating order", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
			return
		}

		for _, p := range products {
			p.OrderID = order.ID
			err := orderRepo.AddOrderProduct(c.Request.Context(), p)
			if err != nil {
				log.Error("Failed adding product to order", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product to order"})
				return
			}
		}

		orderModel.ID = order.ID
		orderModel.UserID = userID.(int64)
		orderModel.CreatedAt = order.CreationDate

		go sendEmail(orderModel)

		c.JSON(http.StatusCreated, gin.H{
			"message": "order created successfully",
			"orderId": order.ID,
		})
	}
}

func GetUserOrdersHandler(orderRepo *orders.Queries, catalogRepo *catalog.Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста
		userID, exists := c.Get("user_id")
		if !exists {
			log.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		orders, err := orderRepo.GetOrders(c.Request.Context(), userID.(int64))
		if err != nil {
			log.Error("Failed getting user orders", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})
			return
		}

		responseOrders := []ResponseOrder{}

		for _, order := range orders {
			products, err := orderRepo.GetOrderProducts(c.Request.Context(), order.ID)
			if err != nil {
				log.Error("Failed getting order", "error", err, "orderID", order.ID)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order"})
				return
			}
			responseOrder := newResponseOrder(catalogRepo, order, products)
			responseOrders = append(responseOrders, responseOrder)
		}

		// Получаем заказы пользователя
		c.JSON(http.StatusOK, gin.H{"orders": responseOrders})
	}
}

func newResponseOrder(catalogRepo *catalog.Repo, order orders.Order, products []orders.OrderProduct) ResponseOrder {
	response := ResponseOrder{
		ID:        order.ID,
		CreatedAt: order.CreationDate.String(),
	}

	for _, orderItem := range products {
		product, err := catalogRepo.GetProduct(orderItem.ProductID)
		item := ResponseOrderProduct{
			Name:     orderItem.ProductName,
			Id:       orderItem.ProductID,
			Quantity: orderItem.Quantity,
			Price:    orderItem.ProductPrice,
		}
		if err != nil {
			log.Error("Failed getting product. It was probably deleted, or added before second DB migration", "error", err, "productID", orderItem.ProductID)
		} else {
			item.ImagePath = product.ImagePath
		}
		response.Products = append(response.Products, item)
	}
	return response
}

type ResponseOrder struct {
	ID        int64                  `json:"id"`
	CreatedAt string                 `json:"createdAt"`
	Products  []ResponseOrderProduct `json:"products"`
}

type ResponseOrderProduct struct {
	Name      string  `json:"productName"`
	Id        string  `json:"productId"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
	ImagePath string  `json:"imagePath"`
}

type Order struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    int64     `json:"userId"`
	Products  []Product `json:"products"`
}
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity int64   `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p *Product) Sum() float64 {
	return float64(p.Quantity) * p.Price
}

func (o *Order) Sum() float64 {
	sum := 0.0
	for _, p := range o.Products {
		sum += p.Sum()
	}
	return sum
}

func (p Product) FormatSum(currency string) string {
	return fmt.Sprintf("%.2f %s", p.Sum(), currency)
}

func (o Order) FormatSum(currency string) string {
	return fmt.Sprintf("%.2f %s", o.Sum(), currency)
}
