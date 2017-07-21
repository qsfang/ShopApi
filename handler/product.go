package handler

import (
	"github.com/labstack/echo"
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
)

//获取商品信息
func GetProInfo(c echo.Context) error {
	var (
		err error
		proid   models.ProductID
		proinfo models.Products
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

	return c.JSON(errcode.ErrSucceed, proinfo)

}

