package models

import (
	"time"

	"ShopApi/orm"
)

type CartsID struct {
	ID			uint64		`json:"id"`
}

type Carts struct {
	ID			uint64		`json:"id"`
	Productid	uint64		`json:"productid"`
	Name		string		`json:"name"`
	Count		uint64		`json:"count"`
	Size		string		`json:"size"`
	Color		string		`json:"color"`
	Imagineid	uint64		`json:"imageid"`
	Userid		uint64		`json:"userid"`
	Status		uint64		`gorm:"column:status" json:"status"`
	Created		time.Time 	`json:"created"`
}

type CartsServiceProvider struct {
}

var CartsService *CartsServiceProvider = &CartsServiceProvider{}

func (cs *CartsServiceProvider) CartsWhether (CartsID uint64)  error {
	var (
		err error
		cart Carts
	)

	db := orm.Conn
	err = db.Where("id = ?", CartsID).First(&cart).Error
	if err != nil {
		return err
	}

	return nil
}

//状态1表示商品在购物车，状态0表示商品不在购物车
func (cs *CartsServiceProvider) CartsDelete (CartsID uint64) error {
	var (
		cart Carts
	)

	db := orm.Conn
	err := db.Model(&cart).Where("id = ?", CartsID).Update("status", 0).Limit(1).Error

	if err != nil {
		return err
	}

	return nil
}
