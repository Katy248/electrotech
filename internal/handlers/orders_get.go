package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"net/http"
)

// GetOrders implements [handler.Handler].
//
//nolint:ireturn
func (h *Handler) GetOrders(ctx context.Context) (handler_gen.GetOrdersRes, error) {
	userID, exists := GetUserID(ctx)
	if !exists {
		return nil, errAuth
	}

	orders, err := h.orderService.GetOrders(ctx, userID)
	if err != nil {
		return &handler_gen.BadResponseStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response: handler_gen.Problem{
				Type:   "get-orders",
				Status: http.StatusUnauthorized,
				Title:  "Failed get orders",
				Detail: handler_gen.OptString{Set: false, Value: ""},
			},
		}, nil
	}

	resp := &handler_gen.GetOrdersResponse{
		Orders: mapOrdersToResponse(orders),
	}

	return resp, nil
}

func mapOrdersToResponse(orders []*models.Order) []handler_gen.Order {
	resp := make([]handler_gen.Order, 0, len(orders))
	for _, o := range orders {
		respOrder := handler_gen.Order{
			ID:           o.ID,
			UserId:       o.ID,
			CreationDate: o.CreationDate.String(),
			Products:     make([]handler_gen.OrderProduct, 0, len(o.OrderProducts)),
		}

		for _, p := range o.OrderProducts {
			respProduct := handler_gen.OrderProduct{
				ID:           p.ID,
				OrderId:      o.ID,
				ProductName:  p.ProductName,
				Quantity:     p.Quantity,
				ProductPrice: p.ProductPrice,
				ProductId:    p.ProductID,
				ImagePath:    p.ImagePath,
			}

			respOrder.Products = append(respOrder.Products, respProduct)
		}

		resp = append(resp, respOrder)
	}

	return resp
}
