package router

import (
	"context"
	"oms_service/orders"
	// "oms_service/orders/requests"

	// "oms_service/database"

	// "oms_service/redis"

	// "net/http"
	"github.com/omniful/go_commons/http"
	// "github.com/omniful/go_commons/jwt/private"
)

func Initialize(ctx context.Context, s *http.Server) (err error) {

	oms_v1 := s.Engine.Group("/api/v1")
	// omsController, ctrlErr := orderwire.Wire(ctx, database.GetClient(), redis.GetClient().Client)
	// 	if ctrlErr != nil {
	// 		err = ctrlErr
	// 		return
	// 	}
	var OrderController *orders.Controller
	// oms_v1.POST("/create-order",omscontroller.CreateOrder)
	// oms_v1.POST("/create-bulk",omscontroller.CreateBulkOrder)
	// oms_v1.GET("/:order_id", private.AuthenticateJWT(), omscontroller.GetOrder)
	// oms_v1.GET("/orders",omscontroller.GetOrders)
	oms_v1.GET("/create", OrderController.CreateOrder)

	return nil

}
