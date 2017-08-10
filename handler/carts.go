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
 *     Initial: 2017/07/21     Zhu Yaqiang
 *     Modify : 2017/07/22     Xu Haosheng
 *     Modify : 2017/07/23     Wang Ke
 *     Modify : 2017/07/24     Ma Chao
 *	   Modify : 2017/08/10     Zhang Zizhao
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

func CartsPutIn(c echo.Context) error {
	var (
		err            error
		carts          models.CartPutIn
		ProInfo        *models.ProductInfo
		prodouctAvatar *string
	)

	if err = c.Bind(&carts); err != nil {
		log.Logger.Error("[ERROR] Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInErrInvalidParams, err.Error())
	}

	if err = c.Validate(&carts); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInErrInvalidParams, err.Error())
	}

	ProInfo, err = models.ProductService.GetProInfo(carts.ProductID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] CartsPutIn: Product doesn't exist", err)

			return general.NewErrorWithMessage(errcode.ErrCartPutInProductNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Get Information with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	prodouctAvatar, err = models.ProductService.GetPrdouctAvatar(carts.ProductID)
	if err != nil {
		log.Logger.Error("[ERROR] Get Information with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMongo, err.Error())
	}

	carts.Avatar = *prodouctAvatar

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.CartsService.CreateInCarts(&carts, userID, ProInfo.Name)
	if err != nil {
		log.Logger.Error("[ERROR] Mysql error with CartPutIn:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInDatabase, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsPutIn name:s%", ProInfo.Name)

	return c.JSON(errcode.ErrCartPutInSucceed, general.NewMessage(errcode.ErrCartPutInSucceed))
}

func CartsDelete(c echo.Context) error {
	var (
		err  error
		cart models.CartDelete
	)

	if err = c.Bind(&cart); err != nil {
		log.Logger.Error("[ERROR] CarDelete Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrCartDeleteErrInvalidParams, err.Error())
	}

	if err = c.Validate(&cart); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCartDeleteErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	UserID := session.Get(general.SessionUserID).(uint64)

	err = models.CartsService.CartsDelete(UserID, cart)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] CarDelete: Product doesn't exist!", err)

			return general.NewErrorWithMessage(errcode.ErrCartDeleteProductNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] CarDelete Mysql:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsDelete ")

	return c.JSON(errcode.ErrCartDeleteSucceed, general.NewMessage(errcode.ErrCartDeleteSucceed))
}

func AlterCartPro(c echo.Context) error {
	var (
		err         error
		cartProduct models.CartAlter
	)

	if err = c.Bind(&cartProduct); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterCartInvalidParams, err.Error())
	}

	if err = c.Validate(cartProduct); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterCartInvalidParams, err.Error())
	}

	err = models.CartsService.AlterCartPro(cartProduct.ID, cartProduct.Count, cartProduct.Color, cartProduct.Size)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] AlterCartspro: Product doesn't exist!", err)

			return general.NewErrorWithMessage(errcode.ErrAlterCartProductNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Alter product with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AlterCartProduct :")

	return c.JSON(errcode.ErrAlterCartSucceed, general.NewMessage(errcode.ErrAlterCartSucceed))
}

func BrowseCart(c echo.Context) error {
	var (
		err    error
		output []models.ConCarts
	)

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	output, err = models.CartsService.BrowseCart(userID)
	if err != nil {
		//if strings.Contains(err.Error(), "not found") {
		//
		//}
		log.Logger.Error("[ERROR] Get Order with error:", err)

		return general.NewErrorWithMessage(errcode.ErrBrowseCartNotFound, err.Error())
	}

	log.Logger.Info("[SUCCEED] BrowseCart")

	return c.JSON(errcode.ErrBrowseCartSucceed, output)
}
