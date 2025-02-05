package orders

import (
	"oms_service/domain"

	"github.com/gin-gonic/gin"
)

// import "oms_service/domain"


type Controller struct {
	// OrderService domain.TenantService
	OrderService domain.OrderService
}

func CreateOrder(c *gin.Context) {
}