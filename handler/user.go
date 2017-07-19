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

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/log"
	"ShopApi/general"
	"ShopApi/orm"
	"ShopApi/general/errcode"
	"ShopApi/server/initorm"
	"ShopApi/models/user"
)


type create struct {
	Mobile 		*string 		`json:"mobile" validate:"required,alphanum,min=6,max=30"`
	Pass 		*string
}

func Create(c echo.Context) error {
	var (
		err 		error
		u 			create
		conn 		orm.Connection
	)

	if err = c.Bind(&u); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	conn, err = initorm.MysqlPool.GetConnection()
	if err != nil {
		log.Logger.Error("Get connection crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrNoConnection, err.Error())
	}
	defer initorm.MysqlPool.ReleaseConnection(conn)

	err = user.UserService.Create(conn, u.Mobile, u.Pass)
	if err != nil {
		log.Logger.Error("create creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
