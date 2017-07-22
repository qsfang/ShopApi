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
 *     Modify: 2017/07/21        Zhang Zizhao 添加创建订单
 *	   Modify: 2017/07/21		 Ai Hao       订单状态更改
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

type Status struct {
	Status uint8 `json:"status"`
}
type ID struct {
	ID uint64 `sql:"auto_increment;primary_key;" json:"id"`
}

type ChangStatus struct {
	ID     uint64 `json:"id"`
	Status uint8  `json:"status"`
}

func CreateOrder(c echo.Context) error {
	var (
		order models.RegisterOrder
		err   error
	)

	if err = c.Bind(&order); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	numberID := sess.Get(general.SessionUserID).(uint64)

	err = models.OrderService.CreateOrder(numberID, order)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("Product not found:", err)

			return general.NewErrorWithMessage(errcode.ErrMysqlfound, err.Error())
		}
		log.Logger.Error("Mysql error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
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
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("Orders not found:", err)

			return general.NewErrorWithMessage(errcode.ErrOrdersNotFound, err.Error())
		}

		log.Logger.Error("Mysql error in get orders:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, orders)
}

func GetOneOrder(c echo.Context) error {
	var (
		err    error
		order  ID
		judge  bool
		OutPut models.GetOrders
	)
	if err = c.Bind(&order); err != nil {
		log.Logger.Error("Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	UserID := sess.Get(general.SessionUserID).(uint64)

	OutPut, err, judge = models.OrderService.GetOneOrder(order.ID, UserID)


	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("Find order with error:", err)

			return general.NewErrorWithMessage(errcode.ErrInformation, err.Error())
		}

		log.Logger.Error("Get Order with error:", err)

		return general.NewErrorWithMessage(errcode.ErrOrdersNotFound, err.Error())
	}

	if judge == false {
		log.Logger.Error("Access with error :", err)

		return general.NewErrorWithMessage(errcode.ErrAccess, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, OutPut)
}

func ChangeStatus(c echo.Context) error {
	var (
		err error
		st  ChangStatus
	)

	if err = c.Bind(&st); err != nil {
		log.Logger.Error("Input order status with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.OrderService.ChangeStatus(st.ID, st.Status)
	if err != nil {
		log.Logger.Error("Change status with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
