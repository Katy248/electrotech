package orders

import (
	"context"
	"electrotech"
	"electrotech/internal/models"
	"electrotech/internal/repository/catalog"
	"electrotech/internal/repository/orders"
	"fmt"
	"net/http"

	"charm.land/log/v2"
	"github.com/gin-gonic/gin"
)

type EmailService interface {
	SendInfo(content []byte, subject string) error
}

type UserRepository interface {
	ByID(ctx context.Context, id int64) (*models.User, error)
}

type Handler struct {
	logger       *log.Logger
	catalogRepo  *catalog.Repo
	ordersRepo   *orders.Repo
	usersRepo    UserRepository
	EmailService EmailService
}

func NewHandler(
	logger *log.Logger,
	catalogRepo *catalog.Repo,
	ordersRepo *orders.Repo,
	usersRepo UserRepository,
	emailService EmailService,
) *Handler {
	return &Handler{
		logger:       logger,
		catalogRepo:  catalogRepo,
		ordersRepo:   ordersRepo,
		usersRepo:    usersRepo,
		EmailService: emailService,
	}
}

type CreateOrderRequest struct {
	Products []OrderProductRequest `binding:"required,min=1" json:"products"`
}

type OrderProductRequest struct {
	ProductId string `binding:"required"       json:"id"`
	Quantity  int    `binding:"required,min=1" json:"quantity"`
}

func (h *Handler) HandleCreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, electrotech.ErrorStr("user not authenticated"))

			return
		}

		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind request")
			c.JSON(http.StatusBadRequest, electrotech.Error(err))

			return
		}

		intUserID, ok := userID.(int64)
		if !ok {
			h.logger.Error("Bad user ID specified", "userID", userID)
			c.JSON(http.StatusNotFound, electrotech.ErrorStr("user not found"))

			return
		}

		user, err := h.usersRepo.ByID(c.Request.Context(), intUserID)
		if err != nil {
			h.logger.Error("User not found", "error", err)
			c.JSON(http.StatusNotFound, electrotech.ErrorStr("user not found"))

			return
		}

		products, err := h.mapProducts(req.Products)
		if err != nil {
			h.logger.Error("Failed map products", "error", err)
			c.JSON(http.StatusNotFound, electrotech.ErrorStr("map products failed"))

			return
		}

		order, err := h.ordersRepo.New(user, products)
		if err != nil {
			h.logger.Error("Failed creating order", "error", err)
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("failed to create order"))

			return
		}

		go h.sendEmail(*order)

		c.JSON(http.StatusCreated, gin.H{
			"orderId": order.ID,
		})
	}
}

func (h *Handler) HandleGetUserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из контекста
		userID, exists := c.Get("user_id")
		if !exists {
			h.logger.Error("User not authenticated")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})

			return
		}

		intUserID, ok := userID.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, electrotech.ErrorStr("bad user ID"))

			return
		}

		orders, err := h.ordersRepo.GetOrders(intUserID)
		if err != nil {
			h.logger.Error("Failed getting user orders", "error", err, "userID", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders"})

			return
		}

		for _, o := range orders {
			for _, p := range o.OrderProducts {
				product, err := h.catalogRepo.GetProduct(p.ProductID)
				if err != nil {
					h.logger.Warn("Failed getting product, skipping", "error", err, "productID", p.ProductID)

					continue
				}

				p.ImagePath = product.ImagePath
			}
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}

func (h *Handler) mapProducts(products []OrderProductRequest) ([]models.OrderProduct, error) {
	mapped := make([]models.OrderProduct, 0, len(products))
	for _, p := range products {
		product, err := h.catalogRepo.GetProduct(p.ProductId)
		if err != nil {
			return nil, fmt.Errorf("product %q: %w", p.ProductId, err)
		}

		mapped = append(mapped,
			models.OrderProduct{ //nolint:exhaustruct_v5
				ProductName:  product.Name,
				Quantity:     int64(p.Quantity),
				ProductPrice: float64(product.Price),
				ProductID:    p.ProductId,
			},
		)
	}

	return mapped, nil
}
