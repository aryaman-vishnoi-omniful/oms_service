package responses

type (
	CreateOrderSvcResponse struct {
		ID              uint64 `json:"id"`
		CustomerID      uint64 `json:"customer_id"`
		TenantID        uint64 `json:"tenant_id"`
		TotalCost       int    `json:"total_cost"`
		ShippingAddress string `json:"shipping_address"`
		BillingAddress  string `json:"billing_address"`
		// Invoice         string `json:"invoice"`
		CurrencyType    string `json:"currency_type"`
		PaymentMethod   string `json:"payment_method"`
		OrderStatus     string `json:"order_status"`
		CreatedBy       uint64 `json:"created_by"`
	}

	CreateOrderCtrlResponse struct {
		CustomerID      uint64 `json:"customer_id"`
		TenantID        uint64 `json:"tenant_id"`
		TotalCost       int    `json:"total_cost"`
		ShippingAddress string `json:"shipping_address"`
		BillingAddress  string `json:"billing_address"`
		// Invoice         string `json:"invoice"`
		CurrencyType    string `json:"currency_type"`
		PaymentMethod   string `json:"payment_method"`
		OrderStatus     string `json:"order_status"`
		CreatedBy       uint64 `json:"created_by"`
	}
)