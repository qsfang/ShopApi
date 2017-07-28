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
 *     Modify : 2017/07/20       Yu Yi
 *     Modify : 2017/07/20       Yang Zhengtian
 *     Modify : 2017/07/27       Li Zebang
 */

package handler

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
)

func AddAddress(c echo.Context) error {
	var (
		err       error
		ormAdress models.OrmAddress
	)

	if err = c.Bind(&ormAdress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrBind, err.Error())
	}

	if err = c.Validate(ormAdress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	ormAdress.UserID = session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.AddAddress(&ormAdress)
	if err != nil {
		log.Logger.Error("[ERROR] AddAddress AddAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func ChangeAddress(c echo.Context) error {
	var (
		err       error
		ormAdress models.OrmAddress
	)

	if err = c.Bind(&ormAdress); err != nil {
		log.Logger.Error("[ERROR] ChangeAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrBind, err.Error())
	}

	if err = c.Validate(ormAdress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.AddressService.FindAddressByAddressID(ormAdress.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.AddressService.ChangeAddress(ormAdress)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeAddress ChangeAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetAddress(c echo.Context) error {
	var (
		err       error
		userId    uint64
		ormAdress models.OrmAddress
		list      []models.AddressGet
	)

	if err = c.Bind(&ormAdress); err != nil {
		log.Logger.Error("[ERROR] GetAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrBind, err.Error())
	}

	if err = c.Validate(ormAdress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	s := session.Get(general.SessionUserID)
	userId = s.(uint64)

	pageStart, pageEnd := utility.Paging(ormAdress.Page, ormAdress.PageSize)
	list, err = models.AddressService.GetAddressByUerID(userId, pageStart, pageEnd)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("Id not find:", err)

			return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
		}
		log.Logger.Error("Mysql err", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, list)
}

func AlterDefault(c echo.Context) error {
	var (
		err error
		m   models.OrmAddress
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrBind, err.Error())
	}

	err = models.AddressService.AlterDefault(m.ID)
	if err != nil {
		log.Logger.Error("Alter Default with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}
	return c.JSON(errcode.ErrSucceed, nil)
}
