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
 *     Modify : 2017/07/22       Xu Haosheng    添加购物车
 *     Modify : 2017/07/23       Wang Ke
 */

package models

import (
	"time"

	"ShopApi/orm"
	"ShopApi/general"
)

type CartsServiceProvider struct {
}

var CartsService *CartsServiceProvider = &CartsServiceProvider{}

type Browse struct {
	Name    string    `json:"name"`
	Count   uint64    `json:"count"`
	Size    string    `json:"size"`
	Color   string    `json:"color"`
	Status  uint64    `json:"status"`
	Created time.Time `json:"created"`
	Type    string    `json:"type"`
	Title   string    `json:"title"`
	Image   string    `json:"image"`
	Url     string    `json:"url"`
}

type Cart struct {
	ProductID uint64    `gorm:"column:productid" json:"productid"`
	ImageID   uint64    `gorm:"column:imageid"json:"imageid"`
	Name      string    `json:"name"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Status    uint64    `json:"status"`
	Created   time.Time `json:"created"`
	Count     uint64    `json:"count"`
}

type Images struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Image string `json:"image"`
	Url   string `json:"url"`
}

type Image struct {
	Type  uint64 `json:"type"`
	Title string `json:"title"`
	Image string `json:"image"`
	Url   string `json:"url"`
	ID    uint64 `json:"id"`
}

func (Image) TableName() string {
	return "image"
}

type CartPro struct {
	ID    uint64 `json:"id"`
	Count uint64 `json:"count"`
	Size  string `json:"size"`
	Color string `json:"color"`
}

type Carts struct {
	ID        uint64    `sql:"primary_key;" gorm:"column:id" json:"id"`
	ProductID uint64    `gorm:"column:productid" json:"productid"`
	Name      string    `json:"name"`
	Count     uint64    `json:"count"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	ImageID   uint64    `gorm:"column:imageid"json:"imageid"`
	Status    uint8     `json:"status"`
	Created   time.Time `json:"created"`
}

type ConCarts struct {
	ID        uint64    `gorm:"column:id" json:"id" validate:"numeric"`
	ProductID uint64    `gorm:"column:productid" json:"productid" validate:"numeric"`
	Name      string    `json:"name" validate:"required, alphaunicode, min = 2, max = 18"`
	Count     uint64    `json:"count" validate:"numeric"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	ImageID   uint64    `gorm:"column:imageid"json:"imageid" validate:"numeric"`
}

// todo:变量
func (cs *CartsServiceProvider) CreateInCarts(carts *ConCarts, userID uint64) error {
	cartsPutIn := Carts{
		UserID:    userID,
		ProductID: carts.ProductID,
		Name:      carts.Name,
		Count:     carts.Count,
		Size:      carts.Size,
		Color:     carts.Color,
		ImageID:   carts.ImageID,
		Created:   time.Now(),
	}

	db := orm.Conn
	err := db.Create(&cartsPutIn).Error

	return err
}

// 状态0表示商品在购物车，状态1表示商品不在购物车
func (cs *CartsServiceProvider) CartsDelete(ID uint64, ProID uint64) error {
	var (
		cart Carts
		err  error
	)

	db := orm.Conn
	err = db.Model(&cart).Where("id = ? AND productid = ?", ID, ProID).Update("status", general.ProductNotInCart).Limit(1).Error

	return err
}

// todo: 返回错误
func (cs *CartsServiceProvider) AlterCartPro(CartsID uint64, Count uint64, Size string, Color string) error {
	var (
		cart Carts
	)
	updater := map[string]interface{}{
		"count": Count,
		"size":  Size,
		"color": Color,
	}

	db := orm.Conn
	err := db.Model(&cart).Where("id = ?", CartsID).Update(updater).Limit(1).Error
	if err != nil {
		return err
	}

	return nil
}

// todo: 命名
//func (cs *CartsServiceProvider) BrowseCart(UserID uint64) ([]Browse, error) {
//	var (
//		err         error
//		carts       []Carts
//		browseCart  []Cart
//		browse      []Browse
//		browseImage Images
//	)
//
//
//	db := orm.Conn
//	err = db.Where("userid = ?", UserID).Find(&carts).Error
//	if err != nil {
//		return browse, err
//	}
//
//	for _, v := range carts {
//		add := Cart {
//			ImageID:  v.ImageID,
//		}
//		browseCart = append(browseCart, add)
//
//		add1 := Browse{
//			Status:  v.Status,
//			Created: v.Created,
//			Count:   v.Count,
//			Name:    v.Name,
//			Color:   v.Color,
//			Size:    v.Size,
//		}
//		browse = append(browse, add1)
//
//		browseImage, err = cs.GetImage(add.ImageID)
//		if err != nil {
//			return browse, err
//		}
//		add2 := Browse{
//			Url:   browseImage.Url,
//			Image: browseImage.Image,
//			Type:  browseImage.Type,
//			Title: browseImage.Title,
//		}
//		browse = append(browse, add2)
//	}
//
//	return browse, err
//}

func (cs *CartsServiceProvider) GetImage(ImageID uint64) (Images, error) {
	var image Images
	db := orm.Conn
	err := db.Where("id = ?", ImageID).Find(&image).Error
	if err != nil {
		return image, err
	}

	return image, nil
}
