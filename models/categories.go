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
 *     Initial: 2017/07/21        Yang Zhengtian
 *     Modify : 2017/07/21        Li Zebang
 */

package models

import (
	"time"

	"ShopApi/orm"
	"ShopApi/general"
)

type CategoriesServiceProvider struct {
}

var CategoriesService *CategoriesServiceProvider = &CategoriesServiceProvider{}

type Categories struct {
	ID      uint64                `sql:"auto_increment;primary_key;",json:"id"`
	Name    string                `json:"name"`
	Pid     uint64                `json:"pid"`
	Status  uint64                `json:"status"`
	Remark  string                `json:"remark"`
	Created time.Time             `json:"created"`
}

type CreateCat struct {
	Name   string `json:"name"`
	Pid    uint64 `json:"pid"`
	Remark string `json:"remark"`
}

func (Categories) TableName() string {
	return "categories"
}

// todo: 代码风格 直接返回错误
func (csp *CategoriesServiceProvider) CheckPid(pid uint64) error {
	var(
		category Categories
	)

	db := orm.Conn

	err:=db.Where("id =? ",pid).First(&category).Error
	if err!=nil {
		return err
	}

	return nil
}

func (csp *CategoriesServiceProvider) Create(ca CreateCat) error {
	cate := Categories{
		Name:                ca.Name,
		Pid:                 ca.Pid,
		Status:              general.CategoriesOnuse,
		Remark:              ca.Remark,
		Created:             time.Now(),
	}

	db := orm.Conn

	err := db.Create(&cate).Error
	if err != nil {
		return err
	}

	return nil
}

// todo: 数据库操作
func (csp *CategoriesServiceProvider) GetCategories(pid uint64) (*[]Categories, error) {
	var (
		categories []Categories
	)

	db := orm.Conn

	return &categories,  db.Where("pid = ? AND status = ?", pid, general.CategoriesOnuse).Find(&categories).Error
}
