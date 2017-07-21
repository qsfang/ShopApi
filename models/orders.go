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
 */

package models

import (
	"time"

	"ShopApi/orm"
	"ShopApi/general"
)

type Orders struct {
	ID         uint64    `sql:"auto_increment;primary_key;" json:"id"`
	UserID     uint64    `gorm:"column:userid" json:"userid"`
	TotalPrice float64   `json:"totalprice"`
	Payment    float64   `json:"payment"`
	Freight    float64   `json:"freight"`
	Remark     string    `json:"remark"`
	Discount   uint8    `json:"discount"`
	Size       string    `json:"size"`
	Color      string    `json:"color"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Payway     uint8     `json:"payway"`
}

type OrderServiceProvider struct {
}

var OrderService *OrderServiceProvider = &OrderServiceProvider{}

func (Orders) TableName() string {
	return "orders"
}

func (osp *OrderServiceProvider) GetOrders(userID uint64, status uint8) ([]Orders, error) {
	var (
		order  Orders
		orders []Orders
	)

	db := orm.Conn

	if status == general.OrderUnfinished || status == general.OrderFinished {
		err := db.Model(&order).Where("userid = ? AND status = ?", userID, status).Find(&orders).Error
		if err != nil {
			return nil, err
		}

		return orders, nil
	}

	err := db.Model(&order).Where("userid = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
