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

	"ShopApi/general"
	"ShopApi/orm"
)

type CategoryServiceProvider struct {
}

var CategoryService *CategoryServiceProvider = &CategoryServiceProvider{}

type Category struct {
	ID      uint64    `sql:"auto_increment;primary_key;",json:"id"`
	Name    string    `json:"name"`
	PID     uint64    `gorm:"column:pid" json:"pid"`
	Status  uint64    `json:"status"`
	Created time.Time `json:"created"`
}

type CreateCategory struct {
	Name string `json:"name" validate:"required,alphanumunicode"`
	PID  uint64 `json:"pid" validate:"required"`
}

type GetCategory struct {
	PID      uint64 `json:"pid"`
	Page     uint64 `json:"page" validate:"required"`
	PageSize uint64 `json:"pagesize" validate:"required"`
}

type CategoryGet struct {
	Name string `json:"name"`
}

func (Category) TableName() string {
	return "category"
}

func (csp *CategoryServiceProvider) CreateCategory(createCategory CreateCategory) error {
	var (
		err error
	)

	category := &Category{
		Name:    createCategory.Name,
		PID:     createCategory.PID,
		Status:  general.CategoryOnUse,
		Created: time.Now(),
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&category).Error

	return err
}

func (csp *CategoryServiceProvider) CheckPID(pid uint64) error {
	var (
		category Category
		err      error
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id =? ", pid).First(&category).Error

	return err
}

func (csp *CategoryServiceProvider) GetCategory(pid, pageStart, pageSize uint64) (*[]CategoryGet, error) {
	var (
		category   Category
		categoryList []CategoryGet
		err        error
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	sql := "SELECT * FROM category WHERE pid = ? AND status = ? LIMIT ?, ? LOCK IN SHARE MODE"

	rows, err := tx.Raw(sql, pid, general.CategoryOnUse, pageStart, pageSize).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tx.ScanRows(rows, &category)
		categoryGet := CategoryGet{Name: category.Name}
		categoryList = append(categoryList, categoryGet)
	}

	return &categoryList, nil
}
