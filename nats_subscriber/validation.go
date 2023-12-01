package sub

import (
	"L0/database"
	"errors"
)

func ValidateOrder(order database.CacheOrder) error {
	if order.Payment.Transaction != order.OrderUID {
		return errors.New("payment does not match")
	}

	for _, item := range order.Items {
		if item.TrackNumber != order.TrackNumber {
			return errors.New("track numbers don't match")
		}
	}

	return nil
}
