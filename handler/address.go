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
 *     Initial: 2017/07/19        Li Zebang
 */

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models/address"
	"ShopApi/orm"
	"ShopApi/server/initorm"
)

type addr struct {
	name      *string `json:"name"`
	phone     *uint64 `json:"phone validate:"required,alphanum,min=6,max=30""`
	province  *string `json:"province"`
	city      *string `json:"city"`
	street    *string `json:"street"`
	address   *string `json:"address"`
	isDefault bool   `json:"is_default"`
}

func Add(c echo.Context) error {
	var (
		err  error
		addr addr
		conn orm.Connection
	)

	if err = c.Bind(&addr); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	conn, err = initorm.MysqlPool.GetConnection()
	if err != nil {
		log.Logger.Error("Get connection crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrNoConnection, err.Error())
	}
	defer initorm.MysqlPool.ReleaseConnection(conn)

	err = address.AddressService.AddAddress(conn, addr.name, addr.province, addr.city, addr.street, addr.address, addr.phone, addr.isDefault)
	if err != nil {
		log.Logger.Error("Add address with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
