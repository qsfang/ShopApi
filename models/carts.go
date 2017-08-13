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
 *     Modify : 2017/08/10       Zhang Zizhao
 */

package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"

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
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	Name      string    `json:"name"`
	Count     uint64    `json:"count"`
	Price     float64   `json:"price"`
	PayStatus uint64    `gorm:"column:paystatus" json:"paystatus"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Status    uint64    `json:"status"`
	Created   time.Time `json:"created"`
}

type CartAlter struct {
	CartsDelete
	ID    uint64 `sql:"primary_key;" validate:"required" json:"id"`
	Count uint64 `json:"count" validate:"required"`
	Color string `json:"color" validate:"required,alphanumunicode"`
	Size  string `json:"size" validate:"required,alphanumunicode"`
}

type ConCarts struct {
	ProductID uint64  `json:"id" validate:"numeric"`
	Name      string  `json:"name" validate:"required,alphaunicode,min=2,max=18"`
	Avatar    string  `json:"avatar" validate:"required"`
	Count     uint64  `json:"num" validate:"numeric"`
	Size      string  `json:"size"`
	Color     string  `json:"color"`
	Price     float64 `json:"price"`
}

type CartsImages struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Avatar    string        `bson:"avatar" json:"avatar"`
}

type CartsSize struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Size      string        `bson:"size" json:"size"`
}

type CartsColors struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Color     string        `bson:"color" json:"color"`
}

type CartPutIn struct {
	ProductID uint64 `json:"productid" validate:"required"`
	Count     uint64 `json:"count" validate:"required"`
	Size      string `json:"size" validate:"required,alphanumunicode"`
	Color     string `json:"color" validate:"required"`
}

type CartDelete struct {
	ProductID uint64 `json:"productid" validate:"required"`
	Size      string `json:"size" validate:"required"`
	Color     string `json:"color" validate:"required"`
}

type CartsDelete struct {
	Data []CartDelete `json:"data"`
}

func (Cart) TableName() string {
	return "cart"
}

func (cs *CartsServiceProvider) CreateCarts(carts *CartPutIn, userID uint64, name string, price float64) error {
	var (
		err        error
		cartsPutIn Cart
		cart       Cart
	)

	db := orm.Conn

	cartsPutIn = Cart{
		UserID:    userID,
		ProductID: carts.ProductID,
		Name:      name,
		Price:     price,
		Count:     carts.Count,
		Size:      carts.Size,
		Color:     carts.Color,
		Status:    general.ProInCart,
		Created:   time.Now(),
	}

	err = db.Where("productid = ? AND userid = ? AND size = ? AND color = ? AND status = ?", carts.ProductID, userID, carts.Size, carts.Color, general.ProInCart).First(&cart).Error
	if err != nil {
		tx := db.Begin()
		defer func() {
			if err != nil {
				err = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}()

		err = tx.Create(&cartsPutIn).Error
		if err != nil {
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			return err
		}

		return err
	}
	count := carts.Count + cart.Count

	err = db.Model(&cart).Where("productid = ? AND size = ? AND color = ? AND status = ?", carts.ProductID, carts.Size, carts.Color, general.ProInCart).Update("count", count).Limit(1).Error

	return err
}

func (cs *CartsServiceProvider) CartDelete(cart *CartDelete, userID uint64) error {
	var (
		ca  Cart
		err error
	)

	ca.ProductID = cart.ProductID
	db := orm.Conn
	err = db.Where("userid = ? AND productid = ? AND size = ? AND color = ? AND status = ?", userID, cart.ProductID, cart.Size, cart.Color, general.ProInCart).Delete(&ca).Error

	return err
}

func (cs *CartsServiceProvider) CartsDelete(data *CartsDelete, userID uint64) error {
	var (
		cart Cart
		err  error
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	for _, delete := range data.Data {
		err = tx.Where("userid = ? AND productid = ? AND size = ? AND color = ? AND status = ?", userID, delete.ProductID, delete.Size, delete.Color, general.ProInCart).Delete(&cart).Error
		if err != nil {
			return err
		}
	}

	return err
}

func (cs *CartsServiceProvider) AlterCartPro(CartsID uint64, Count uint64) error {
	var (
		cart Cart
		err  error
	)
	db := orm.Conn
	err = db.Where("id = ?", CartsID).Find(&cart).Error
	if err != nil {
		return err
	}

	updater := map[string]interface{}{"count": Count}
	err = db.Model(&cart).Where("id = ?", CartsID).Update(updater).Limit(1).Error

	return err
}

func (cs *CartsServiceProvider) CartsBrowse(userID uint64) (*[]ConCarts, error) {
	var (
		err   error
		cart  []Cart
		image ProductImages
		list  []ConCarts
	)

	db := orm.Conn

	err = db.Where("status = ? AND userid = ?", general.ProInCart, userID).Find(&cart).Error
	if err != nil {
		return nil, err
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for _, value := range cart {
		err = collection.Find(bson.M{"productid": value.ProductID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return nil, err
		}
		lis := ConCarts{
			ProductID: value.ProductID,
			Name:      value.Name,
			Color:     value.Color,
			Count:     value.Count,
			Size:      value.Size,
			Price:     value.Price,
			Avatar:    image.Image,
		}
		list = append(list, lis)
	}

	return &list, err
}
