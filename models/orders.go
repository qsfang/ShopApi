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
 *	   Modify : 2017/07/21		 Ai Hao       订单状态更改
 *	   Modify : 2017/07/21		 Zhang Zizhao 创建订单
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

type Order struct {
	ID         uint64    `sql:"auto_increment;primary_key;"json:"id"`
	UserID     uint64    `gorm:"column:userid" json:"userid"`
	AddressID  uint64    `gorm:"column:addressid" json:"addressid"`
	TotalPrice float64   `gorm:"column:totalprice" json:"totalprice"`
	Freight    float64   `gorm:"column:freight" json:"freight"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	Status     uint8     `gorm:"column:status" json:"status"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	PayWay     uint8     `gorm:"column:payway" json:"payway"`
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
	ImageID    uint64    `gorm:"column:imageid"json:"imageid"`
	Remark     string    `json:"remark"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

func (Order) TableName() string {
	return "order"
}

func (osp *OrderServiceProvider) CreateOrder(numberID uint64, o OrmOrders) error {
	var (
		err error
		car Cart
	)

	order := Order{
		UserID:     numberID,
		AddressID:  o.AddressID,
		TotalPrice: o.TotalPrice,
		Freight:    o.Freight,
		Remark:     o.Remark,
		Status:     general.OrderFinished,
		PayWay:     o.PayWay,
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
		"status":    0,
		"paystatus": 0,
		"orderid":   order.ID,
	}
	err = tx.Model(&car).Where("userid = ? AND status = 1 AND paystatus = 1", numberID).Update(changeMap).Limit(1).Error

	return err
}

func (osp *OrderServiceProvider) GetOrders(userID uint64, status uint8, pageStart, pageEnd uint64) (*[]Order, error) {
	var (
		order  Order
		orders []Order
	)

	db := orm.Conn

	if status == general.OrderUnfinished || status == general.OrderFinished {
		sql := fmt.Sprintf("SELECT * FROM order WHERE userid = ? AND status = ? LIMIT %d, %d LOCK IN SHARE MODE", pageStart, pageEnd)

		rows, err := db.Raw(sql, userID, status).Rows()
		defer rows.Close()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			db.ScanRows(rows, &order)
			orders = append(orders, order)
		}

		return &orders, nil
	}

	sql := fmt.Sprintf("SELECT * FROM order WHERE userid = ? LIMIT %d, %d LOCK IN SHARE MODE", pageStart, pageEnd)

	rows, err := db.Raw(sql, userID).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		db.ScanRows(rows, &order)
		orders = append(orders, order)
	}

	return &orders, nil
}

func (osp *OrderServiceProvider) GetOneOrder(UserID uint64, ID uint64) ([]OrmOrders, error) {
	var (
		err      error
		order    Order
		carts    []Cart
		getOrder []OrmOrders
	)

	db1 := orm.Conn
	err = db1.Where("id = ? AND userid = ?", ID, UserID).First(&order).Error
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

	db2 := orm.Conn
	err = db2.Where("orderid = ?", order.ID).Find(&carts).Error
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
	cha := Order{
		Status: status,
	}

	updater := map[string]interface{}{"status": status}
	db := orm.Conn

	err := db.Model(&cha).Where("id=?", id).Update(updater).Limit(1).Error
	return err
}
