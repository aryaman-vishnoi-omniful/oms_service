package domain

import (
	"context"
	// "oms_service/orders"
	"oms_service/orders/requests"
	// "github.com/omniful/go_commons/file_utilities/requests"
)

type OrderService interface{
	CreateTenant(c context.Context, request *requests.CreateOrderSvcRequest)
}