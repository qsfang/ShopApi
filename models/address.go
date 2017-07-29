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
 *     Initial: 2017/07/18        Li Zebang
 *     Modify : 2017/07/20        Yu Yi
 *     Modify : 2017/07/20        Yang Zhengtian
 *     Modify : 2017/07/28        Li Zebang
 */

package models

import (
	"time"

	"ShopApi/general"
	"ShopApi/orm"
)

type AddressServiceProvider struct {
}

var AddressService *AddressServiceProvider = &AddressServiceProvider{}

type Address struct {
	ID        uint64    `sql:"auto_increment; primary_key;" json:"id"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Street    string    `json:"street"`
	Address   string    `json:"address"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	IsDefault uint8     `gorm:"column:isdefault" json:"isdefault" `
}

type AddAddress struct {
	UserID    uint64 `json:"userid"`
	Name      string `json:"name" validate:"required,alphanumunicode"`
	Phone     string `json:"phone" validate:"required,numeric,len=11"`
	Province  string `json:"province" validate:"required,alphanumunicode"`
	City      string `json:"city" validate:"required,alphanumunicode"`
	Street    string `json:"street" validate:"required,alphanumunicode"`
	Address   string `json:"address" validate:"required,alphanumunicode"`
	IsDefault uint8  `json:"isdefault" validate:"max=1,min=0"`
}

type ChangeAddress struct {
	ID       uint64 `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,alphanumunicode"`
	Phone    string `json:"phone" validate:"required,numeric,len=11"`
	Province string `json:"province" validate:"required,alphanumunicode"`
	City     string `json:"city" validate:"required,alphanumunicode"`
	Street   string `json:"street" validate:"required,alphanumunicode"`
	Address  string `json:"address" validate:"required,alphanumunicode"`
}

type GetAddress struct {
	UserID   uint64 `json:"userid"`
	Page     uint64 `json:"page" validate:"required"`
	PageSize uint64 `json:"pagesize" validate:"required"`
}

type AddressGet struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Street   string `json:"street"`
	Address  string `json:"address"`
}

type AddressAlter struct {
	ID     uint64 `json:"id" validate:"required"`
	UserID uint64 `json:"userid"`
}

func (Address) TableName() string {
	return "address"
}

func (csp *AddressServiceProvider) AddAddress(addAddress *AddAddress) (err error) {
	var (
		address *Address
	)

	address = &Address{
		UserID:    addAddress.UserID,
		Name:      addAddress.Name,
		Phone:     addAddress.Phone,
		Province:  addAddress.Province,
		City:      addAddress.City,
		Street:    addAddress.Street,
		Address:   addAddress.Address,
		Created:   time.Now(),
		Updated:   time.Now(),
		IsDefault: addAddress.IsDefault,
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(address).Error

	return err
}

func (csp *AddressServiceProvider) ChangeAddress(changeAddress *ChangeAddress) (err error) {
	var (
		address Address
	)

	updater := map[string]interface{}{
		"name":     changeAddress.Name,
		"phone":    changeAddress.Phone,
		"province": changeAddress.Province,
		"city":     changeAddress.City,
		"street":   changeAddress.Street,
		"address":  changeAddress.Address,
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Model(&address).Where("id = ?", changeAddress.ID).Update(updater).Limit(1).Error

	return err
}

func (csp *AddressServiceProvider) FindAddressByAddressID(ID uint64) (err error) {
	var (
		address Address
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ?", ID).First(&address).Error

	return err
}

func (csp *AddressServiceProvider) GetAddressByUserID(userID uint64, pageStart, pageSize uint64) (addressList *[]AddressGet, err error) {
	var (
		address     Address
		addresses  []AddressGet
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	sql := "SELECT * FROM address WHERE userid = ? LIMIT ?, ? LOCK IN SHARE MODE"

	rows, err := tx.Raw(sql, userID, pageStart, pageSize).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tx.ScanRows(rows, &address)
		addressGet := AddressGet{
			Province: address.Province,
			City:     address.City,
			Street:   address.Street,
			Address:  address.Address,
		}
		addresses = append(addresses, addressGet)
	}

	return &addresses, nil
}

func (csp *AddressServiceProvider) AlterAddressToDefault(alterAddress *AddressAlter) (err error) {
	var (
		address Address
	)

	updater := map[string]interface{}{"isdefault": general.AddressDefault}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Model(&address).Where("id = ?", alterAddress.ID).Update(updater).Limit(1).Error

	return err
}

func (csp *AddressServiceProvider) AlterAddressToNotDefault(userID uint64) (err error) {
	var (
		address Address
	)

	updater := map[string]uint8{"isdefault": general.AddressNotDefault}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	return orm.Conn.Model(&address).Where("userid = ?", userID).Update(updater).Limit(1).Error
}
