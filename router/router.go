package router

import (
	"context"
	"oms_service/orders"

	// "oms_service/orders/requests"

	// "oms_service/database"

	// "oms_service/redis"

	// "net/http"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/config"
	// "github.com/omniful/go_commons/jwt/private"
)

func Initialize(ctx context.Context, s *http.Server) (err error) {
	s.Engine.Use(log.RequestLogMiddleware(log.MiddlewareOptions{
		Format:      config.GetString(ctx, "log.format"),
		Level:       config.GetString(ctx, "log.level"),
		LogRequest:  config.GetBool(ctx, "log.request"),
		LogResponse: config.GetBool(ctx, "log.response"),
	}))

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
	oms_v1.POST("/create", OrderController.CreateOrder)
	// oms_v1.POST("/create-bulkorder", OrderController.CreateBulkCsv)

	return nil

}
