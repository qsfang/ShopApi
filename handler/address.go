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
	"errors"

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
		err        error
		addAddress models.AddressJSON
	)

	if err = c.Bind(&addAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAddAddressInvalidParams, err.Error())
	}

	log.Logger.Info("%#v", addAddress)

	if err = c.Validate(addAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAddAddressInvalidParams, err.Error())
	}

	addAddress.UserID = utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request()).Get(general.SessionUserID).(uint64)

	err = models.AddressService.AddAddress(&addAddress)
	if err != nil {
		log.Logger.Error("[ERROR] AddAddress AddAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AddAddress: UserID %d", addAddress.UserID)

	return c.JSON(errcode.ErrSucceed, general.NewMessage(errcode.AddAddressSucceed))
}

func ChangeAddress(c echo.Context) error {
	var (
		err           error
		changeAddress models.AddressJSON
	)

	if err = c.Bind(&changeAddress); err != nil {
		log.Logger.Error("[ERROR] ChangeAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeAddressInvalidParams, err.Error())
	}

	if err = c.Validate(changeAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeAddressInvalidParams, err.Error())
	}

	err = models.AddressService.FindAddressByAddressID(changeAddress.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrChangeAddressNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.AddressService.ChangeAddress(&changeAddress)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeAddress ChangeAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangeAddress: UserID %d", changeAddress.ID)

	return c.JSON(errcode.ChangeAddressSucceed, general.NewMessage(errcode.ChangeAddressSucceed))
}

func GetAddress(c echo.Context) error {
	var (
		err         error
		userID      uint64
		addressList *[]models.AddAddress
	)

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID = session.Get(general.SessionUserID).(uint64)

	addressList, err = models.AddressService.GetAddressByUserID(userID)
	if err != nil {
		log.Logger.Error("[ERROR] GetAddress GetAddressByUserID ", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if len(*addressList) == 0 {
		err = errors.New("[ERROR] Address Not Found")

		log.Logger.Error("[ERROR] GetAddress GetAddressByUserID:", err)

		return general.NewErrorWithMessage(errcode.ErrGetAddressNotFound, err.Error())
	}

	log.Logger.Info("[SUCCEED] Get address by userId: %d", userID)

	return c.JSON(errcode.ErrSucceed, addressList)
}

func AlterDefault(c echo.Context) error {
	var (
		err          error
		alterAddress *models.AlterAddress
	)

	if err = c.Bind(&alterAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterDefaultInvalidParams, err.Error())
	}

	if err = c.Validate(alterAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterDefaultInvalidParams, err.Error())
	}

	err = models.AddressService.FindAddressByAddressID(alterAddress.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] AlterDefault FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrAlterDefaultNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] AlterDefault FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	alterAddress.UserID = session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.AlterAddress(alterAddress)
	if err != nil {
		log.Logger.Error("[ERROR] AlterDefault AlterAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AlterDefault by userId:", alterAddress.UserID)

	return c.JSON(errcode.AlterDefaultSucceed, general.NewMessage(errcode.AlterDefaultSucceed))
}
