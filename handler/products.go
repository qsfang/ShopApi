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
 *     Initial: 2017/07/21        Ai Hao
 *     Modify: 2017/07/21         Yu Yi
 */

package handler

import (
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"github.com/labstack/echo"
)

//名称name，totalsale  ，类型categories，价格price，原价originalprice，
// 状态status，尺码siez，颜色color,封面图片imageid，图片集imageids，评论remark,
//详细信息 detail ，创建日期 created，存货量inventory

func CreateP(c echo.Context) error {
	var (
		err error
		p   models.CreatePro
	)

	if err = c.Bind(&p); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.CreateP(p)
	if err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetProductList(c echo.Context) error {
	var (
		err    	error
		m       models.GetCategories
		list 	[]models.GetProList
	)

	if err = c.Bind(&m); err != nil {
	log.Logger.Error("Get categories with error:", err)

	return general.NewErrorWithMessage(errcode.ErrMysql,err.Error())
	}

	list, err = models.ProductService.GetProduct(m)
	if err != nil {
	log.Logger.Error("Error", err)

	return general.NewErrorWithMessage(errcode.ErrMysql,err.Error())
	}

	return c.JSON(errcode.ErrSucceed, list)
}

func ChangeProStatus(c echo.Context) error {
	var(
		err		error
		pro		models.ChangePro
	)

	if err = c.Bind(&pro); err != nil {
		log.Logger.Error("Change crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeProStatus(pro)
	if err != nil {
		log.Logger.Error("change crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

//根据商品ID获取商品信息
func GetProInfo(c echo.Context) error {
	var (
		err error
		proid   models.ProductID
		proinfo models.Product
	)

	if err = c.Bind(&proid); err != nil {
		log.Logger.Error("Get crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	proinfo,err = models.ProductService.GetProInfo(proid)

	if err != nil {
		log.Logger.Error("error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Debug("i got here :%v", proinfo)

	return c.JSON(errcode.ErrSucceed, proinfo)
}
