package orders

import (
	"log"
	// "net/http"
	"oms_service/domain"
	"oms_service/orders/requests"

	"github.com/gin-gonic/gin"
	// commonError "github.com/omniful/go_commons/error"
)

// import "oms_service/domain"


type Controller struct {
	// OrderService domain.TenantService
	OrderService domain.OrderService
}

func(tc *Controller) CreateOrder(c *gin.Context) {

	var createOrderReq *requests.CreateOrderCtrlRequest
	err:=c.ShouldBind(&createOrderReq)
	if err!=nil{
		// cuserr:=commonError.NewCustomError("Bad Request",err.Error())
		log.Fatal("bad request")
		return


	}


}