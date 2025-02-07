package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// import "time"

type CreateOrderSvcRequest struct {
	CustomerId      uint64
	TenantId        uint64
	TotalCost       int
	ShippingAddress string
	BillingAddress  string
	// Invoice string
	CurrencyType  string
	PaymentMethod string
	OrderStatus   string
	// created_at time.Time
	// updated_at time.Time
}

type CreateOrderCtrlRequest struct {
	CustomerID      string         `json:"customer_id"`
	TenantID        string         `json:"tenant_id" `
	TotalCost       int            `json:"total_cost" `
	ShippingAddress string         `json:"shipping_address" `
	BillingAddress  string         `json:"billing_address"`
	Invoice         map[string]any `json:"invoice"`
	CurrencyType    string         `json:"currency_type" `
	PaymentMethod   string         `json:"payment_method"`
	OrderStatus     string         `json:"order_status" `
}

type OrderItem struct {
	OrderID  string`bson:"order_id"`
	SKUID    string `bson:"sku_id,omitempty" json:"sku_id" `
	Quantity int    `bson:"quantity" json:"quantity"   `
}

// Order represents the main order entity
type Order struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id" `
	CustomerName string              `bson:"customer_name,omitempty" json:"customer_name"`
	OrderNo      string              `bson:"order_no" json:"order_no"`
	OrderItems   []OrderItem         `bson:"order_items" json:"order_items"`
	Status       string              `bson:"status" json:"status"`
	CreatedAt    primitive.DateTime  `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    primitive.DateTime  `bson:"updated_at,omitempty" json:"updated_at"`
	DeletedAt    *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at" `
}
