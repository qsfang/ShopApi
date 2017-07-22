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
// todo:对齐
/*
 * Revision History:
 *     Initial: 2017/07/21        Zhu Yaqiang
 *     Modify: 2017/07/22     Xu Haosheng    添加购物车
 */

package models

import (
	"time"

	"ShopApi/orm"
)

type CartsServiceProvider struct {
}

var CartsService *CartsServiceProvider = &CartsServiceProvider{}

type CartsDel struct {
	ID    uint64 `gorm:"column:id" json:"id"`
	ProID uint64 `json:"productid"`
}

type Test struct {
	ID  uint64 `gorm:"column:id" json:"id"`
	UserID uint64  `json:"userid"`
}

type Browse struct {
	ProductID uint64    `gorm:"column:productid" json:"productid"`
	Name      string    `json:"name"`
	Count     uint64    `json:"count"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	ImageID   uint64    `gorm:"column:imageid"json:"imageid"`
	Status    uint64    `json:"status"`
	Created   time.Time `json:"created"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Url       string    `json:"url"`
}
type Cart struct {
	ProductID uint64    `gorm:"column:productid" json:"productid"`
	ImageID   uint64    `gorm:"column:imageid"json:"imageid"`
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
	Type  string `json:"type"`
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
	Status    uint64    `json:"status"`
	Created   time.Time `json:"created"`
}

// todo:变量
func (cs *CartsServiceProvider) CreateInCarts(carts Carts, userID uint64) error {
	cartsPutIn := Carts{
		UserID:    userID,
		ProductID: carts.ProductID,
		Name:      carts.Name,
		Count:     carts.Count,
		Size:      carts.Size,
		Color:     carts.Color,
		ImageID:   carts.ImageID,
		Status:    carts.Status,
		Created:   time.Now(),
	}

	db := orm.Conn

	err := db.Create(&cartsPutIn).Error
	if err != nil {
		return err
	}

	return nil
}

// 状态0表示商品在购物车，状态1表示商品不在购物车
// todo: 常量定义  数据库操作
func (cs *CartsServiceProvider) CartsDelete(ID uint64, ProID uint64) error {
	var (
		cart Carts
		err  error
	)

	db := orm.Conn

	err = db.Where("id = ? and productid = ?", ID, ProID).First(&cart).Error
	if err != nil {
		return err
	}

	err = db.Model(&cart).Where("id = ? and productid = ?", ID, ProID).Update("status", 1).Limit(1).Error
	if err != nil {
		return err
	}

	return nil
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
func (cs *CartsServiceProvider) BrowseCart(UserID uint64) ([]Browse, error) {
	var (
		err         error
		carts       Carts
		browse      []Browse
		browsepro   []Product
		browseimage []Images
	)

	db := orm.Conn
	err = db.Where("userid = ?", UserID).Find(&carts).Error
	if err != nil {
		return browse, err
	}
	var cart = Cart{
		ProductID: carts.ProductID,
		ImageID:   carts.ImageID,
		Status:    carts.Status,
		Created:   carts.Created,
		Count:     carts.Count,
	}

	browsepro, err = cs.GetProduct(cart.ProductID)
	if err != nil {
		return browse, err
	}

	for _, x := range browsepro {
		add := Browse{
			Name:  x.Name,
			Size:  x.Size,
			Color: x.Color,
		}
		browse = append(browse, add)
	}

	browseimage, err = cs.GetImage(cart.ImageID)
	if err != nil {
		return browse, err
	}

	for _, s := range browseimage {
		add := Browse{
			Url:   s.Url,
			Image: s.Image,
			Type:  s.Type,
			Title: s.Title,
		}
		browse = append(browse, add)
	}

	return browse, err
}

func (cs *CartsServiceProvider) GetProduct(ProductID uint64) ([]Product, error) {
	var product []Product
	db := orm.Conn
	err := db.Where("id = ?", ProductID).Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (cs *CartsServiceProvider) GetImage(ImageID uint64) ([]Images, error) {
	var image []Images
	db := orm.Conn
	err := db.Where("id = ?", ImageID).Find(&image).Error
	if err != nil {
		return image, err
	}

	return image, nil
}
