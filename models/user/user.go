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

package user

import (
	"time"

	"github.com/jinzhu/gorm"

	"ShopApi/orm"
	"ShopApi/utility"
	"ShopApi/general"
)

type UserServiceProvider struct{
}

var UserService *UserServiceProvider = &UserServiceProvider{}

type User struct {
	UserID 		uint64		`orm:"pk;column(id)"  json:"userid"`
	OpenID 		string		`orm:"column(openid)" json:"openid"`
	Name 		string		`json:"name"`
	Pass 		string		`json:"pass"`
	Status		uint16		`json:"status"`
	Type		uint16		`json:"type"`
	Created 	time.Time	`json:"created"`
}

func (User) TableName() string {
	return "user"
}

func (us *UserServiceProvider) Create(conn orm.Connection, name, pass *string) error {
	hashedPass, err := utility.GenerateHash(*pass)
	if err != nil {
		return err
	}

	u := User{
		Name: 		*name,
		Pass:		string(hashedPass),
		Status:		general.UserActive,
		Type: 		general.PhoneUser,
		Created:	time.Now(),
	}

	db := conn.(*gorm.DB)

	err = db.Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}
