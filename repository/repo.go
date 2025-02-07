package repository

import (
	"context"
	"oms_service/database"
	"oms_service/orders/requests"
	// "oms_service/orders/responses"
	"time"
	// "github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/exchange_rate/responses"
)
type OrderService struct{
	
}

func CreateOrder(c context.Context,order *requests.Order)(error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Database("OMS_servicee").Collection("orders")

	_, err := collection.InsertOne(ctx, order)
	return err



}