/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/07/21       Li Zebang
 *	   Modify : 2017/07/21		 Ai Hao
 *	   Modify : 2017/07/21		 Zhang Zizhao
 *     Modify : 2017/07/21       Ma Chao
 */

package models

import (
	"fmt"
	"time"

	"ShopApi/general"
	"ShopApi/orm"
)

type OrderServiceProvider struct {
}

var OrderService *OrderServiceProvider = &OrderServiceProvider{}

type Orders struct {
	ID         uint64    `sql:"auto_increment;primary_key;"json:"id"`
	UserID     uint64    `gorm:"column:userid" json:"userid"`
	AddressID  uint64    `gorm:"column:addressid" json:"addressid"`
	TotalPrice float64   `gorm:"column:totalprice" json:"totalprice"`
	PayWay     uint8     `gorm:"column:payway" json:"payway"`
	Freight    float64   `json:"freight"`
	Remark     string    `json:"remark"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type OrmOrders struct {
	ID         uint64    `json:"id" validate:"required,numeric"`
	UserID     uint64    `json:"userid" validate:"required,numeric"`
	TotalPrice float64   `json:"totalprice" validate:"required,numeric"`
	Payment    float64   `json:"payment" validate:"required,numeric"`
	Freight    float64   `json:"freight" validate:"required,numeric"`
	Discount   uint8     `json:"discount" validate:"numeric"`
	Name       string    `json:"name"validate:"required, alphaunicode, min = 2, max = 18"`
	Size       string    `json:"size" validate:"required,alphanum"`
	Count      uint64    `json:"count"validate:"required,numeric"`
	Color      string    `json:"color" validate:"required,alphanum"`
	Status     uint8     `json:"status" validate:"required,numeric"`
	PayWay     uint8     `json:"payway" validate:"required,numeric"`
	Page       uint64    `json:"page" validate:"required,numeric"`
	PageSize   uint64    `json:"pagesize" validate:"required,numeric"`
	AddressID  uint64    `json:"addressid" validate:"required,numeric"`
	ImageID    uint64    `json:"imageid"`
	Remark     string    `json:"remark"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type CreateOrder struct {
	AddressID  uint64  `json:"addressid" validate:"required"`
	TotalPrice float64 `json:"totalprice" validate:"required"`
	Freight    float64 `json:"freight" validate:"required"`
	Payment    float64 `json:"payment" validate:"required"`
	Remark     string  `json:"remark"`
	PayWay     uint8   `json:"payway" validate:"required"`
}

type GetOrders struct {
	UserID   uint64 `json:"userid"`
	Status   uint8  `json:"status" validate:"required,max=239,min=236"`
	Page     uint64 `json:"page" validate:"required"`
	PageSize uint64 `json:"pagesize" validate:"required"`
}

type OrdersGet struct {
	TotalPrice float64 `json:"totalprice"`
	Freight    float64 `json:"freight"`
	Remark     string  `json:"remark"`
	Status     uint8   `json:"status"`
}

func (Orders) TableName() string {
	return "orders"
}

func (osp *OrderServiceProvider) CreateOrder(numberID uint64, ord CreateOrder) error {
	var (
		err error
		car Cart
	)

	order := Orders{
		UserID:     numberID,
		AddressID:  ord.AddressID,
		TotalPrice: ord.TotalPrice,
		Freight:    ord.Freight,
		Remark:     ord.Remark,
		Status:     general.OrderUnfinished,
		PayWay:     ord.PayWay,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	db := orm.Conn

	tx := db.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&order).Error
	if err != nil {
		return err
	}

	changeMap := map[string]interface{}{
		"status":  general.ProNotInCart,
		"orderid": order.ID,
	}
	err = tx.Model(&car).Where("userid = ? AND status = ? AND paystatus = ?", numberID, general.ProInCart, general.Buy).Update(changeMap).Limit(1).Error

	return err
}

func (osp *OrderServiceProvider) GetOrders(getOrders *GetOrders, pageStart uint64) (*[]OrdersGet, error) {
	var (
		sql        string
		order      Orders
		orderGet   OrdersGet
		ordersList []OrdersGet
	)

	db := orm.Conn

	if getOrders.Status == general.OrderUnfinished || getOrders.Status == general.OrderFinished || getOrders.Status == general.OrderCanceled {
		sql = fmt.Sprintf("SELECT * FROM orders WHERE userid = %d AND status = %d LIMIT %d, %d LOCK IN SHARE MODE", getOrders.UserID, getOrders.Status, pageStart, getOrders.PageSize)
	} else {
		sql = fmt.Sprintf("SELECT * FROM orders WHERE userid = %d LIMIT %d, %d LOCK IN SHARE MODE", getOrders.UserID, pageStart, getOrders.PageSize)
	}

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		db.ScanRows(rows, &order)
		fmt.Println(order)
		orderGet = OrdersGet{
			TotalPrice: order.TotalPrice,
			Freight:    order.Freight,
			Remark:     order.Remark,
			Status:     order.Status,
		}
		ordersList = append(ordersList, orderGet)
	}

	return &ordersList, nil
}

func (osp *OrderServiceProvider) GetOneOrder(UserID uint64, ID uint64) ([]OrmOrders, error) {
	var (
		err      error
		order    Orders
		carts    []Cart
		getOrder []OrmOrders
	)

	db := orm.Conn
	tx := db.Begin()

	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ? AND userid = ?", ID, UserID).First(&order).Error
	if err != nil {
		return getOrder, err
	}

	add1 := OrmOrders{
		TotalPrice: order.TotalPrice,
		Freight:    order.Freight,
		Status:     order.Status,
		Created:    order.Created,
		PayWay:     order.PayWay,
		Updated:    order.Updated,
		AddressID:  order.AddressID,
	}
	getOrder = append(getOrder, add1)

	err = tx.Where("orderid = ?", order.ID).Find(&carts).Error
	if err != nil {
		return getOrder, err
	}

	for _, v := range carts {
		add1 := OrmOrders{
			Name:    v.Name,
			Count:   v.Count,
			Size:    v.Size,
			Color:   v.Color,
			ImageID: v.ImageID,
		}
		getOrder = append(getOrder, add1)
	}

	return getOrder, nil
}

func (osp *OrderServiceProvider) ChangeStatus(id uint64, status uint8) error {
	cha := Orders{
		Status: status,
	}

	updater := map[string]interface{}{"status": status}
	db := orm.Conn

	err := db.Model(&cha).Where("id = ?", id).Update(updater).Limit(1).Error
	return err
}
