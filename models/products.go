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
 *     Initial: 2017/07/21         Ai Hao
 *     Modify : 2017/07/21         Zhu Yaqiang
 *     Modify : 2017/07/21         Yu Yi
 *     Modify : 2017/07/21         Machao
 */

package models

import (
	"fmt"
	"time"

	"ShopApi/general"
	"ShopApi/orm"
	"gopkg.in/mgo.v2/bson"
)

type ProductServiceProvider struct {
}

var ProductService *ProductServiceProvider = &ProductServiceProvider{}

type Product struct {
	ID        uint64    `sql:"auto_increment;primary_key;" gorm:"column:id" json:"id"`
	Name      string    `json:"name"`
	TotalSale uint64    `gorm:"column:totalsale" json:"totalsale"`
	Category  uint64    `json:"categories"`
	Price     float64   `json:"price"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	Detail    string    `json:"detail"`
	Status    uint8     `json:"status"`
	Created   time.Time `json:"created"`
}

type ProductImages struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Class     uint8         `bson:"class" json:"class"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Image     string        `bson:"image" json:"image"`
}

type CreateProduct struct {
	Name         string   `json:"name" validate:"required"`
	Avatar       string   `json:"avatar"`
	Images       []string `json:"images"`
	DetailImages []string `json:"detailimages"`
	Category     uint64   `json:"category" validate:"required"`
	Price        float64  `json:"price" validate:"required"`
	Size         string   `json:"size" validate:"required,alphanumunicode"`
	Color        string   `json:"color" validate:"required,alphanumunicode"`
	Detail       string   `json:"detail" validate:"required"`
}

type ProductList struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Avatar   string  `json:"avatar"`
	Category uint64  `json:"categories"`
	Price    float64 `json:"price"`
}

type ProductCategory struct {
	Category uint64 `json:"category"`
	Page     uint64 `json:"page" validate:"required, numeric"`
	PageSize uint64 `json:"pagesize" validate:"required, numeric"`
}

type ProductID struct {
	ID uint64 `json:"id" validate:"required"`
}

type ProductInfo struct {
	Name         string   `json:"name"`
	Images       []string `json:"images"`
	DetailImages []string `json:"detailimages"`
	TotalSale    uint64   `json:"totalsale"`
	Category     uint64   `json:"category"`
	Price        float64  `json:"price"`
	Size         string   `json:"size"`
	Color        string   `json:"color"`
	Detail       string   `json:"detail"`
}

type ChangeProStatus struct {
	ID     uint64 `json:"id" validate:"required"`
	Status uint8  `json:"status" validate:"required"`
}

type ChangeCategory struct {
	ID       uint64 `json:"id" validate:"required"`
	Category uint64 `json:"category" validate:"required"`
}

func (Product) TableName() string {
	return "product"
}

func (ps *ProductServiceProvider) CreateProduct(create *CreateProduct) error {
	var (
		err     error
		product Product
	)

	product = Product{
		Name:     create.Name,
		Category: create.Category,
		Price:    create.Price,
		Size:     create.Size,
		Color:    create.Color,
		Detail:   create.Detail,
		Status:   general.ProductOnSale,
		Created:  time.Now(),
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&product).Error
	if err != nil {
		return err
	}

	err = tx.Where("name = ?", create.Name).First(&product).Error
	if err != nil {
		return err
	}

	err = AddProductImage(product.ID, create)

	return err
}

func AddProductImage(productID uint64, create *CreateProduct) error {
	var (
		err    error
		image  ProductImages
		images []ProductImages
	)

	image = ProductImages{
		Class:     general.ProductAvatar,
		ProductID: productID,
		Image:     create.Avatar,
	}

	images = append(images, image)

	for _, img := range create.Images {
		image = ProductImages{
			Class:     general.ProductImage,
			ProductID: productID,
			Image:     img,
		}
		images = append(images, image)
	}

	for _, img := range create.DetailImages {
		image = ProductImages{
			Class:     general.ProductDetailImage,
			ProductID: productID,
			Image:     img,
		}
		images = append(images, image)
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for _, img := range images {
		err = collection.Insert(img)
		if err != nil {
			break
		}
	}

	return err
}

func (ps *ProductServiceProvider) GetProductHeader() (*[]string, error) {
	var (
		err     error
		product Product
		image   ProductImages
		header  []string
	)

	db := orm.Conn
	
	sql := "SELECT * FROM product WHERE status = ? LIMIT 5 LOCK IN SHARE MODE"
	rows, err := db.Raw(sql, general.ProductOnSale).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductImage}).One(&image)
		if err != nil {
			return nil, err
		}
		header = append(header, image.Image)
	}

	return &header, nil
}

func (ps *ProductServiceProvider) GetProductList() (*[]ProductList, error) {
	var (
		err     error
		product ProductList
		image   ProductImages
		list    []ProductList
	)

	db := orm.Conn

	sql := "SELECT * FROM product WHERE status = ? LIMIT 5, 6 LOCK IN SHARE MODE"
	rows, err := db.Raw(sql, general.ProductOnSale).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return nil, err
		}
		product.Avatar = image.Image
		list = append(list, product)
	}

	return &list, nil
}

func (ps *ProductServiceProvider) GetProductByCategory(cate, pageStart, pageSize uint64) (*[]ProductList, error) {
	var (
		product ProductList
		image   ProductImages
		list    []ProductList
	)

	db := orm.Conn

	sql := fmt.Sprintf("SELECT * FROM product WHERE category = ? AND status = ? LIMIT %d, %d LOCK IN SHARE MODE", pageStart, pageSize)

	rows, err := db.Raw(sql, cate, general.ProductOnSale).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return nil, err
		}
		product.Avatar = image.Image
		list = append(list, product)
	}

	return &list, nil
}

func (ps *ProductServiceProvider) GetProInfo(id uint64) (*ProductInfo, error) {
	var (
		err     error
		product Product
		info    ProductInfo
		images  []ProductImages
	)

	db := orm.Conn

	err = db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}

	info = ProductInfo{
		Name:      product.Name,
		TotalSale: product.TotalSale,
		Category:  product.Category,
		Price:     product.Price,
		Size:      product.Size,
		Color:     product.Color,
		Detail:    product.Detail,
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	err = collection.Find(bson.M{"productid": id}).All(&images)
	if err != nil {
		return nil, err
	}

	fmt.Println(images)

	for _, image := range images {
		switch image.Class {
		case general.ProductAvatar:
			continue
		case general.ProductImage:
			info.Images = append(info.Images, image.Image)
		case general.ProductDetailImage:
			info.DetailImages = append(info.DetailImages, image.Image)
		}
	}

	return &info, nil
}

func (ps *ProductServiceProvider) ChangeProStatus(sta *ChangeProStatus) error {
	var (
		pro Product
	)

	updater := map[string]uint8{"status": sta.Status}

	db := orm.Conn

	return db.Model(&pro).Where("id = ?", sta.ID).Updates(updater).Limit(1).Error
}

func (ps *ProductServiceProvider) ChangeCategory(cate *ChangeCategory) error {
	var (
		pro Product
	)

	db := orm.Conn
	err := db.Model(&pro).Where("id = ?", cate.ID).Update("category", cate.Category).Limit(1).Error

	return err
}
