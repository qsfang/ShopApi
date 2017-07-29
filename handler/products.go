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
 *      Initial: 2017/07/21         Ai Hao
 *      Modify : 2017/07/21         Zhu Yaqiang
 *      Modify : 2017/07/21         Yu Yi
 *      Modify : 2017/07/21         Ma Chao
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

func CreateProduct(c echo.Context) error {
	var (
		err       error
		product   models.CreateProduct
	)

	if err = c.Bind(&product); err != nil {
		log.Logger.Error("[ERROR] CreateProduct Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if err = c.Validate(product); err != nil {
		log.Logger.Error("[ERROR] CreateProduct Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.CreateProduct(&product)
	if err != nil {
		log.Logger.Error("[ERROR] CreateProduct CreateProduct:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetProductList(c echo.Context) error {
	var (
		err  error
		cate models.ConProduct
		list *[]models.ConProduct
	)

	if err = c.Bind(&cate); err != nil {
		log.Logger.Error("[ERROR] GetProductList Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	pageStart := utility.Paging(cate.Page, cate.PageSize)
	list, err = models.ProductService.GetProduct(cate.Category, pageStart, cate.PageSize)
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] Categories not exist", err)

			return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
		}

		log.Logger.Error("[ERROR] Get categories with error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, list)
}

func ChangeProStatus(c echo.Context) error {
	var (
		err error
		pro models.ChangeProStatus
	)

	if err = c.Bind(&pro); err != nil {
		log.Logger.Error("[ERROR] Change Product Status Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if err = c.Validate(pro); err != nil {
		log.Logger.Error("[ERROR] Product Status Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if pro.Status != general.ProductOnsale && pro.Status != general.ProductUnsale {
		err = errors.New("[ERROR] Product Status InExistent")
		log.Logger.Error("[ERROR] status transformed with error :",err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeProStatus(&pro)
	if err != nil {
		log.Logger.Error("[ERROR] status transformed with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetProInfo(c echo.Context) error {
	var (
		err           error
		ProInfo       *models.ConProduct
		ProInfoReturn *models.Product
	)

	if err = c.Bind(&ProInfo); err != nil {
		log.Logger.Error("[ERROR] GetProInfo Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	ProInfoReturn, err = models.ProductService.GetProInfo(ProInfo.ID)

	if err != nil {
		log.Logger.Error("[ERROR] Get info with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, ProInfoReturn)
}

func ChangeCategories(c echo.Context) error {
	var (
		err error
		m   *models.ChangeCategories
	)

	if err = c.Bind(&m); err != nil {
		log.Logger.Error("[ERROR] ChangeCategories Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if err = c.Validate(m); err != nil {
		log.Logger.Error("[ERROR] Categories Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeCategories(m)
	if err != nil {

		log.Logger.Error("[ERROR] Categories change with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
