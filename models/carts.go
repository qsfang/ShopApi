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
 *     Initial: 2017/07/21       Zhu Yaqiang
 *     Modify : 2017/07/22       Xu Haosheng
 *     Modify : 2017/07/23       Wang Ke
 *     Modify : 2017/07/24       Ma Chao
 */

package models

import (
	"time"

	"ShopApi/general"
	"ShopApi/orm"
)

type CartsServiceProvider struct {
}

var CartsService *CartsServiceProvider = &CartsServiceProvider{}

type Cart struct {
	ID        uint64    `sql:"primary_key;" gorm:"column:id" json:"id"`
	ProductID uint64    `gorm:"column:productid" json:"productid"`
	OrderID   uint64    `gorm:"column:orderid" json:"orderid"`
	Count     uint64    `json:"count" validate:"required,numeric"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	ImageID   uint64    `gorm:"column:imageid" json:"imageid"`
	PayStatus uint64    `gorm:"column:paystatus" json:"paystatus"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Name      string    `json:"name"`
	Status    uint64    `json:"status"`
	Created   time.Time `json:"created"`
}

type ConCarts struct {
	ProductID uint64    `json:"productid" validate:"numeric"`
	Status    uint64    `json:"status" validate:"required, numeric, max = 1"`
	Name      string    `json:"name" validate:"required, alphaunicode, min = 2, max = 18"`
	Count     uint64    `json:"count" validate:"numeric"`
	ImageID   uint64    `json:"imageid" validate:"numeric"`
	Created   time.Time `json:"created"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
}

type CartUse struct {
	Name      string    `json:"name"`
	ProductID uint64    `json:"productid"`
	Count     uint64    `json:"count"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	ImageID   uint64    `json:"imageid"`
	Created   time.Time `json:"created"`
}

type CartDel struct {
	ProductID uint64    `json:"productid"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
}

func(Cart) TableName() string {
	return "cart"
}

func (cs *CartsServiceProvider) CreateInCarts(carts *CartUse, userID uint64) error {
	var (
		err                     error
		cartsPutIn         Cart
	)

	cartsPutIn = Cart{
		UserID:                   userID,
		ProductID:             carts.ProductID,
		Name:                    carts.Name,
		Count:                    carts.Count,
		Size:                        carts.Size,
		Color:                     carts.Color,
		ImageID:                carts.ImageID,
		Created:                 time.Now(),
		Status:                    general.ProInCart,
	}

	db := orm.Conn
	err = db.Create(&cartsPutIn).Error

	return err
}

func (cs *CartsServiceProvider) CartsDelete(UserID uint64, carts *CartDel) error {
	var (
		cart Cart
		err  error
	)

	db := orm.Conn
	err = db.Model(&cart).Where("productid = ? AND userid = ? AND size = ? AND color = ?", carts.ProductID, UserID, carts.Size, carts.Color).Update("status", general.ProNotInCart).Limit(1).Error

	return err
}


func (cs *CartsServiceProvider) AlterCartPro(CartsID uint64, Count uint64, Color string, Size string) error {
	var (
		cart Cart
	)

	updater := map[string]interface{}{"count": Count, "color": Color, "size": Size }

	db := orm.Conn
	err := db.Model(&cart).Where("id = ?", CartsID).Update(updater).Limit(1).Error

	return err
}

func (cs *CartsServiceProvider) BrowseCart(UserID uint64) ([]ConCarts, error) {
	var (
		err         error
		carts       []Cart
		browse      []ConCarts
	)

	db := orm.Conn
	err = db.Where("userid = ? AND status = ?", UserID, general.ProInCart).Find(&carts).Error
	if err != nil {
		return browse, err
	}

	for _, v := range carts {
		add1 := ConCarts{
			ImageID:   v.ImageID,
			Status:    v.Status,
			Created:   v.Created,
			Count:     v.Count,
			Name:      v.Name,
			Color:     v.Color,
			Size:      v.Size,
			ProductID: v.ProductID,
		}
		browse = append(browse, add1)
	}

	return browse, err
}
