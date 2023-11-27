package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := "host=localhost user=natsdbuser password=deeznats dbname=natsdb port=5432 sslmode=disable TimeZone=GMT"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetFromDB(orderUID string, db *gorm.DB) (*CacheOrder, error) {
	order := DBOrder{}
	payment := Payment{}
	delivery := Delivery{}
	items := []OrderItem{}

	result := db.First(&order, "order_uid = ?", orderUID)
	if result.Error != nil {
		return nil, result.Error
	}
	db.First(&delivery, "order_uid = ?", order.OrderUID)
	db.First(&payment, "transaction = ?", order.OrderUID)
	db.Find(&items, "track_number = ?", order.TrackNumber)
	return &CacheOrder{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
	}, nil
}

func SaveToDB(order CacheOrder, db *gorm.DB) error {
	dborder := DBOrder{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
	delivery := order.Delivery
	payment := order.Payment
	items := order.Items

	result := db.Create(&dborder)
	if result.Error != nil {
		return result.Error
	}

	result = db.Create(&delivery)
	if result.Error != nil {
		return result.Error
	}

	result = db.Create(&payment)
	if result.Error != nil {
		return result.Error
	}

	for _, item := range items {
		result = db.Create(&item)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
