package repository

import (
	"context"
	"oms_service/orders/requests"
	"oms_service/orders/responses"
	// "github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/exchange_rate/responses"
)
type OrderService interface {
	CreateOrder(c context.Context, request *requests.CreateOrderSvcRequest) (*responses.CreateOrderSvcResponse,error)
}


func CreateOrder(c context.Context,svcReq *requests.CreateOrderSvcRequest){


}