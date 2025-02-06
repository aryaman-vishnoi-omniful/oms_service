package requests

// import "time"

type CreateOrderSvcRequest struct{
	CustomerId uint64
	TenantId uint64
	TotalCost int
	ShippingAddress string
	BillingAddress string
	// Invoice string
	CurrencyType string
	PaymentMethod string
	OrderStatus string
	// created_at time.Time
	// updated_at time.Time
}


type CreateOrderCtrlRequest struct {
	CustomerID      string            `json:"customer_id"`
	TenantID        string            `json:"tenant_id" `
	TotalCost       int               `json:"total_cost" `
	ShippingAddress string            `json:"shipping_address" `
	BillingAddress  string            `json:"billing_address"`
	Invoice         map[string]any    `json:"invoice"`
	CurrencyType    string            `json:"currency_type" `
	PaymentMethod   string            `json:"payment_method"`
	OrderStatus     string            `json:"order_status" `
}
