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
 *		Initial: 2017/07/21			Ai Hao
 *		Modify: 2017/07/21			Zhu Yaqiang
 *      Modify: 2017/07/21          Yu Yi
 */

package handler

import (
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
)

func CreateProduct(c echo.Context) error {
	var (
		err error
		p   models.CreatePro
	)

	if err = c.Bind(&p); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.CreateProduct(p)
	if err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

// todo：代码规范
func GetProductList(c echo.Context) error {
	var (
		err  error
		m    models.GetCategories
		list []models.GetProList
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("Get categories with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	list, err = models.ProductService.GetProduct(m)
	if err != nil {
		log.Logger.Error("Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, list)
}

func ChangeProStatus(c echo.Context) error {
	var (
		err error
		pro models.ChangePro
	)

	if err = c.Bind(&pro); err != nil {
		log.Logger.Error("Change crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeProStatus(pro)
	if err != nil {
		log.Logger.Error("change chanslates with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

// 根据商品ID获取商品信息
func GetProInfo(c echo.Context) error {
	var (
		err     error
		ProID   models.ProductID
		ProInfo models.Product
	)

	if err = c.Bind(&ProID); err != nil {
		log.Logger.Error("Analysis crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	ProInfo, err = models.ProductService.GetProInfo(ProID.ID)

	if err != nil {
		log.Logger.Error("Get info with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, ProInfo)
}

func ChangeCategories(c echo.Context) error {
	var (
		err error
		m   models.ChangeCate
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("Categories change with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeCategories(m)
	if err != nil {
		log.Logger.Error("Categories change with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
