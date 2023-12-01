package database

import (
	"time"
)

//OrderUID == Payment.Transaction

//TrackNumber == OrderItem.TrackNumber

type CacheOrder struct {
	OrderUID          string      `json:"order_uid" fake:"-"`
	TrackNumber       string      `json:"track_number" fake:"-"`
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

type DBOrder struct {
	OrderUID          string `gorm:"primaryKey" fake:"-"`
	TrackNumber       string `fake:"-"`
	Entry             string
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	Shardkey          string
	SmID              int
	DateCreated       time.Time
	OofShard          string
}

type Delivery struct {
	OrderUID string `gorm:"primaryKey" fake:"-" json:",omitempty"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" fake:"-"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type OrderItem struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number" fake:"-" `
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
