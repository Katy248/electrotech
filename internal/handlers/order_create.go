package handlers

import (
	"context"
	"electrotech/internal/gen/handler_gen"
	"electrotech/internal/models"
	"net/http"
)

// CreateOrder implements [handler.Handler].
//
//nolint:ireturn
func (h *Handler) CreateOrder(
	ctx context.Context,
	req *handler_gen.CreateOrderReq,
) (handler_gen.CreateOrderRes, error) {
	userID, ok := GetUserID(ctx)
	if !ok {
		return nil, errAuth
	}

	items := mapProducts(req.Products)

	orderID, err := h.orderService.CreateOrder(ctx, userID, items)
	if err != nil {
		h.logger.Error("Failed creating order", "error", err)

		return &handler_gen.BadResponseStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: handler_gen.Problem{
				Type:   "create-order",
				Status: http.StatusInternalServerError,
				Title:  "Failed create order",
				Detail: handler_gen.OptString{Set: false, Value: ""},
			},
		}, nil
	}

	return &handler_gen.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

func mapProducts(products []handler_gen.CreateOrderReqProductsItem) []models.OrderItem {
	mapped := make([]models.OrderItem, 0, len(products))
	for _, p := range products {
		mapped = append(mapped,
			models.OrderItem{
				Quantity:  p.Quantity,
				ProductID: p.ID,
			},
		)
	}

	return mapped
}
