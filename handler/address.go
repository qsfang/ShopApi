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
 *     Initial: 2017/07/19       Li Zebang
 *     Modify: 2017/07/20        Yu Yi
 *     Modify: 2017/07/20        Yang Zhengtian
 */

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
)

type Add struct {
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	Province  *string `json:"province"`
	City      *string `json:"city"`
	Street    *string `json:"street"`
	Address   *string `json:"address"`
	IsDefault int8    `json:"isdefault"`
}


func AddAddress(c echo.Context) error {
	var (
		err  error
		addr Add
	)

	if err = c.Bind(&addr); err != nil {
		log.Logger.Error("Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	contact := &models.Contact{
		UserID:    userID,
		Name:      *addr.Name,
		Phone:     *addr.Phone,
		Province:  *addr.Province,
		City:      *addr.City,
		Street:    *addr.Street,
		Address:   *addr.Address,
		IsDefault: addr.IsDefault,
	}

	err = models.ContactService.AddAddress(contact)
	if err != nil {
		log.Logger.Error("Add address with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}


func ChangeAddress(c echo.Context) error {
	var (
		err error
		m   Address
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ContactService.ChangeAddress(m.ID, m.Name, m.Phone, m.Province, m.City, m.Street, m.Address)
	if err != nil {
		log.Logger.Error("create creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
>>>>>>> 1210e443c562e0452a2a5ad5315891f74cb1b621

func GetAddress(c echo.Context) error {
	var (
		err    error
		userid uint64
		list   []models.Addressget
	)

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	s := sess.Get(general.SessionUserID)
	userid = s.(uint64)
	list, err = models.ContactService.GetAddress(userid)
	if err != nil {
		log.Logger.Error("error:", err)
		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, list)
}

func Alter(c echo.Context) error {
	var (
		err error
	)
	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	user := session.Get(general.SessionUserID).(uint64)

	err = models.ContactService.AlterDefalt(user)
	if err != nil {
		log.Logger.Error("Alter Default with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}
	return c.JSON(errcode.ErrSucceed, nil)
}