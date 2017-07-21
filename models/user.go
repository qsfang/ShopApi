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
 *	   Modify: 2017/07/20		  Zhang Zizhao  登录检查
 */

package models

import (
	"time"

	"ShopApi/general"
	"ShopApi/orm"
	"ShopApi/utility"
)

type UserServiceProvider struct {
}

var UserService *UserServiceProvider = &UserServiceProvider{}

type User struct {
	UserID   uint64    `sql:"auto_increment;primary_key;" gorm:"column:id" json:"userid"`
	OpenID   string    `gorm:"column:openid" json:"openid"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Status   uint16    `json:"status"`
	Type     uint16    `json:"type"`
	Created  time.Time `json:"created"`
}

type UserInfo struct {
	UserID   uint64 `sql:"primary_key" gorm:"column:userid" json:"userid"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Sex      uint8  `json:"sex"`
}


func (User) TableName() string {
	return "users"
}

func (UserInfo) TableName() string {
	return "userinfo"
}

func (us *UserServiceProvider) Create(name, pass *string) error {
	hashedPass, err := utility.GenerateHash(*pass)
	if err != nil {
		return err
	}

	u := User{
		Name:     *name,
		Password: string(hashedPass),
		Status:   general.UserActive,
		Type:     general.PhoneUser,
		Created:  time.Now(),
	}

	db := orm.Conn

	tx := db.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&u).Error
	if err != nil {
		return err
	}

	info := UserInfo{
		UserID: u.UserID,
		Phone:  *name,
		Sex:    general.Man,
	}

	err = tx.Create(&info).Error
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceProvider) Login(name, pass *string) (bool, uint64, error) {
	var (
		u   User
		err error
	)

	db := orm.Conn
	err = db.Where("name = ?", *name).First(&u).Error
	if err!=nil {

		return false, 0, err
	}

	if !utility.CompareHash([]byte(u.Password), *pass)  {

		return false, 0, nil
	}

	return true, u.UserID, nil
}

func (us *UserServiceProvider) GetInfo(UserID uint64) (UserInfo, error) {

	var (
		err error
		s   UserInfo
	)

	db := orm.Conn
	err = db.Where("UserID = ?", UserID).First(&s).Error
	if err != nil {
		return s, err
	}

	return s, nil
}

func (us *UserServiceProvider)ChangePhone(UserID uint64,Phone *string) error{
	var (
		err	error
		con	Contact
	)
	change := map[string]interface{}{"phone": Phone}

	db := orm.Conn
	err = db.Model(&con).Where("userid=?",UserID).Update(change).Limit(1).Error
	if err != nil{
		return err
	}
	return  nil
}

func (us *UserServiceProvider) ChangeMobilePassword(oldpass *string ,newpass *string , id uint64) (bool,error) {
	var(
		user   User
		err error
	)
	db := orm.Conn
	err = db.Where("id = ?", id).First(&user).Error
	if err!=nil {
		return false,err
	}

	if !utility.CompareHash([]byte(user.Password), *oldpass)  {

		return false, nil
	}
	err = db.Model(&user).Where("id = ?", id).Update(user.Password,*newpass).Limit(1).Error
	if err!=nil{
		return false,err
	}

	return true, nil
}
