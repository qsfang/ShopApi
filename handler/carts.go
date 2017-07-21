package handler

import (
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
)

func Cartsdel(c echo.Context) error {
	var (
		err error
		cartid   models.CartsID
	)

	if err = c.Bind(&cartid); err != nil {
		log.Logger.Error("Get crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.CartsService.CartsWhether(cartid.ID)
	if err == gorm.ErrRecordNotFound {
		log.Logger.Error("The product doesn't exist !", err)

		return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
	}

	err = models.CartsService.CartsDelete(cartid.ID)

	if err != nil {
		log.Logger.Error("Delete product with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
