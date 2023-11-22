package database

import (
	"time"

	"github.com/google/uuid"
)

//OrderUID == Payment.Transaction

//TrackNumber == OrderItem.TrackNumber

type Order struct {
	OrderUID          uuid.UUID   `json:"order_uid" gorm:"primaryKey" fake:"-"`
	TrackNumber       uuid.UUID   `json:"track_number" fake:"-"`
	Entry             string      `json:"entry"`
	Delivery          Delivery    `json:"delivery" fake:"-"`
	Payment           Payment     `json:"payment" fake:"-"`
	Items             []OrderItem `json:"items" fake:"-"`
	Locale            string      `json:"locale"`
	InternalSignature string      `json:"internal_signature"`
	CustomerID        string      `json:"customer_id"`
	DeliveryService   string      `json:"delivery_service"`
	Shardkey          string      `json:"shardkey"`
	SmID              int         `json:"sm_id"`
	DateCreated       time.Time   `json:"date_created"`
	OofShard          string      `json:"oof_shard"`
}

type Delivery struct {
	OrderUID uuid.UUID `gorm:"primaryKey;foreignKey:Order.OrderUID" fake:"-"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Zip      string    `json:"zip"`
	City     string    `json:"city"`
	Address  string    `json:"address"`
	Region   string    `json:"region"`
	Email    string    `json:"email"`
}

type Payment struct {
	Transaction  uuid.UUID `json:"transaction" gorm:"primaryKey" fake:"-"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDt    int       `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type OrderItem struct {
	ChrtID      int       `json:"chrt_id" gorm:"primaryKey;autoIncrement"`
	TrackNumber uuid.UUID `json:"track_number" fake:"-"`
	Price       int       `json:"price"`
	Rid         string    `json:"rid"`
	Name        string    `json:"name"`
	Sale        int       `json:"sale"`
	Size        string    `json:"size"`
	TotalPrice  int       `json:"total_price"`
	NmID        int       `json:"nm_id"`
	Brand       string    `json:"brand"`
	Status      int       `json:"status"`
}
