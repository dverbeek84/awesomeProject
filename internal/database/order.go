package database

import (
	"awesomeProject/internal/model"
)

func MigrateOrderScheme() error {
	return DB.AutoMigrate(
		model.Order{},
	)
}

func DummyOrderData() error {
	return nil
}

func CreateOrder(order *model.Order) error {
	if err := DB.Create(order).Error; err != nil {
		return err
	}
	return nil
}

// UpdateOrderDoneByID is used to update an order record with done.
func UpdateOrderDoneByID(id string) error {
	if err := DB.Model(&model.Order{}).Where("id = ?", id).Update("state", "done").Error; err != nil {
		return err
	}
	return nil
}
