package requests

type CreateOrderSvcRequest struct{
	TotalCost int
	ShippingAddress string
	BillingAddress string
	Invoice string
	CurrencyType string
	PaymentMethod string
	OrderStatus string
}