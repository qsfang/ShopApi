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
 *     Initial: 2017/07/18        Yusan Kurban
 */

package address

import (
	"time"

	"ShopApi/orm"
	"ShopApi/log"
)

type Address struct {
	ID        uint64    `sql:"auto_increment;primary_key;" json:"id"`
	Name      string    `json:"name"`
	Phone     uint64    `json:"phone"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Address   string    `json:"address"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	Created   time.Time `json:"created"`
	IsDefault bool      `gorm:"column:isdefault" json:"isdefault"`
}

type AddressServiceProvider struct {
}

var AddressService *AddressServiceProvider = &AddressServiceProvider{}

func (as *AddressServiceProvider) AddAddress(name *string, province *string, city *string, street *string, address *string, phone *uint64, isDefault bool) error {
	log.Logger.Debug("name :%s, province :%s, city :%s, street :%s, address :%s, phone :%d, isDefault :%v", *name, *province, *city, *street, *address, *phone, isDefault)
	addr := &Address{
		Name:      *name,
		Phone:     *phone,
		Province:  *province,
		City:      *city,
		Street:    *street,
		Address:   *address,
		UserID:    1,
		Created:   time.Now(),
		IsDefault: isDefault,
	}

	db := orm.Conn

	err := db.Create(&addr).Error
	if err != nil {
		return err
	}

	return nil
}
