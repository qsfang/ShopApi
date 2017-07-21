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
 *     Initial: 2017/07/21       Li Zebang
 *     Modify: 2017/07/21        Zhang Zizhao //添加创建订单
 */

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
	"github.com/jinzhu/gorm"
)

type Status struct {
	Status uint8 `json:"status"`
}


func CreateOrder(c echo.Context) error {
	var (
		order models.Registerorder
		err   error
	)

	if err = c.Bind(&order); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}
	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	numberID := sess.Get(general.SessionUserID).(uint64)

	err = models.OrderService.Createorder(numberID, order)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("User not found:", err)

			return general.NewErrorWithMessage(errcode.ErrNamefound, err.Error())
		}
		if err == gorm.ErrInvalidTransaction {
			log.Logger.Error("no valid transaction", err)

			return general.NewErrorWithMessage(errcode.ErrNamefound, err.Error())
		} else {
			log.Logger.Error("Mysql error:", err)

			return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
		}
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetOrders(c echo.Context) error {
	var (
		err    error
		status Status
		orders []models.Orders
	)

	if err = c.Bind(&status); err != nil {
		log.Logger.Error("Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if status.Status != general.OrderUnfinished && status.Status != general.OrderFinished && status.Status != general.OrderGetAll {
		return general.NewErrorWithMessage(errcode.ErrInvalidOrdersStatus, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	orders, err = models.OrderService.GetOrders(userID, status.Status)
	if err != nil {
		log.Logger.Error("Get orders with error:", err)
		return general.NewErrorWithMessage(errcode.ErrGetOrders, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, orders)
}
