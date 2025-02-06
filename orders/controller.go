package orders

import (
	"log"
	"strconv"

	// "net/http"
	// "oms_service/domain"
	"oms_service/orders/requests"
	"oms_service/orders/responses"
	"oms_service/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	error2 "github.com/omniful/go_commons/error"
	oresponse "github.com/omniful/go_commons/response"
	// commonError "github.com/omniful/go_commons/error"
)

// import "oms_service/domain"


type Controller struct {
	// OrderService domain.TenantService
	OrderService repository.OrderService
}

func(tc *Controller) CreateOrder(c *gin.Context) {

	var createOrderReq *requests.CreateOrderCtrlRequest
	err:=c.ShouldBind(&createOrderReq)
	if err!=nil{
		// cuserr:=commonError.NewCustomError("Bad Request",err.Error())
		log.Fatal("bad request")
		return


	}
	svcRequest, err := convertControllerReqToServiceReqCreateOrder(c, createOrderReq)
	if err!=nil{
		// tenantError.NewErrorResponse(c, cusErr
		log.Fatal(err.Error())
		return
	}

	svcResponse, err := tc.OrderService.CreateOrder(c,svcRequest)
	if err!=nil{
		// tenantError.NewErrorResponse(c, cusErr)
		log.Fatal(err.Error())
		return
	}

	response := convertServiceRespToControllerRespCreateOrder(svcResponse)
	oresponse.NewSuccessResponse(c, response)


}

func convertControllerReqToServiceReqCreateOrder(ctx *gin.Context, createOrderReq *requests.CreateOrderCtrlRequest) (svcReq *requests.CreateOrderSvcRequest, cusErr error2.CustomError) {
	validate := validator.New()
	err := validate.Struct(createOrderReq)
	if err != nil {
		return nil, error2.NewCustomError("400", "validation failed")
	}

	tenantID, err := strconv.ParseUint(createOrderReq.TenantID, 10, 64)
	if err != nil {
		return nil, error2.NewCustomError("PARSE_INT_ERROR", err.Error())
	}

	customerID, err := strconv.ParseUint(createOrderReq.CustomerID, 10, 64)
	if err != nil {
		return nil, error2.NewCustomError("PARSE_INT_ERROR", err.Error())
	}
	if createOrderReq.OrderStatus==""{
		createOrderReq.OrderStatus="on-hold"
	}

	svcReq = &requests.CreateOrderSvcRequest{
		CustomerId:      customerID,
		TenantId:        tenantID,
		TotalCost:       createOrderReq.TotalCost,
		ShippingAddress: createOrderReq.ShippingAddress,
		BillingAddress:  createOrderReq.BillingAddress,
		// Invoice:         createOrderReq.Invoice,
		CurrencyType:    createOrderReq.CurrencyType,
		PaymentMethod:   createOrderReq.PaymentMethod,
		OrderStatus:     createOrderReq.OrderStatus,
	}

	return 
}



func convertServiceRespToControllerRespCreateOrder(resp *responses.CreateOrderSvcResponse) *responses.CreateOrderCtrlResponse {
	return &responses.CreateOrderCtrlResponse{
		CustomerID:      resp.CustomerID,
		TenantID:        resp.TenantID,
		TotalCost:       resp.TotalCost,
		ShippingAddress: resp.ShippingAddress,
		BillingAddress:  resp.BillingAddress,
		// Invoice:         resp.Invoice,
		CurrencyType:    resp.CurrencyType,
		PaymentMethod:   resp.PaymentMethod,
		OrderStatus:     resp.OrderStatus,
		CreatedBy:       resp.CreatedBy,
	}
}
