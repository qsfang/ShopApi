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
 *     Initial: 2017/07/19        Yu yi, Li Zebang
 */

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
)

type Address struct {
<<<<<<< HEAD
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	Province  *string `json:"province"`
	City      *string `json:"city"`
	Street    *string `json:"street"`
	Address   *string `json:"address"`
	IsDefault int8    `json:"isdefault"`
=======
	ID        *uint64  `sql:"auto_increment; primary_key;" json:"id"`
	Name      *string  `json:"name"`
	Phone     *string  `json:"phone"`
	Province  *string  `json:"province"`
	City      *string  `json:"city"`
	Street    *string  `json:"street"`
	Address   *string  `json:"address"`
	IsDefault bool     `json:"isdefault"`
>>>>>>> a96214f93b05e8b310e4ad6265f8d4db8e33e3b0
}

func Add(c echo.Context) error {
	var (
		err  error
		addr Address
	)

<<<<<<< HEAD
	if err = c.Bind(&addr); err != nil {
=======

	if err = c.Bind(&m); err != nil {
>>>>>>> a96214f93b05e8b310e4ad6265f8d4db8e33e3b0
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

<<<<<<< HEAD
	//session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	//user := session.Get(general.SessionUserID).(uint64)
	user := uint64(166)

	err = models.ContactService.AddAddress(addr.Name, &user, addr.Phone, addr.Province, addr.City, addr.Street, addr.Address, addr.IsDefault)
=======
	err = models.ContactService.ChangeAddress(m.ID, m.Name, m.Phone, m.Province, m.City, m.Street, m.Address)
>>>>>>> a96214f93b05e8b310e4ad6265f8d4db8e33e3b0
	if err != nil {
		log.Logger.Error("Add address with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func ChangeAddress(c echo.Context) error {
	var (
		err   error
		m     Address
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ContactService.ChangeAddress(m.Name, m.Province, m.City, m.Street, m.Address)
	if err != nil {
		log.Logger.Error("create creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
