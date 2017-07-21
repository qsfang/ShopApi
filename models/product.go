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
 *     Initial: 2017/07/20        Zhu Yaqiang
 */
package models

import (
	"ShopApi/orm"
	"time"
)


type ProductServiceProvider struct {
}

var ProductService *ProductServiceProvider = &ProductServiceProvider{}


type ProductID struct{
	ID				uint64 `json:"id"`
}

type Products struct {
	ID				uint64 `sql:"primary_key" gorm:"column:id" json:"id"`
	Name			string `json:"name"`
	Totalsale   	uint64 `json:"totalsale"`
	Categories		uint64 `json:"categories"`
	Price			float64 `json:"price"`
	Originalprice	float64 `json:"originalprice"`
	Status          uint64 `json:"status"`
	Size            string `json:"size"`
	Color           string `json:"color"`
	Imageid			uint64 `json:"imageid"`
	Imageids		string `json:"imageids"`
	Remark			string `json:"remark"`
	Detail			string `json:"detail"`
	Created			time.Time `json:"created"`
	Inventory		uint64 `json:"inventory"`
}

func (proinfoser *ProductServiceProvider) GetProInfo(ProID ProductID) (Products,error) {

	var (
		err error
		proinfo   Products
	)

	db := orm.Conn
	err = db.Where("id = ?", ProID.ID).First(&proinfo).Error
	if err != nil {
		return proinfo, err
	}
	return proinfo, nil
}