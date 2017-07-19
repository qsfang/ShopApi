package handler

import (
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/orm"
	"ShopApi/utility"

	"github.com/labstack/echo"
)

func LoginHandlerMobilephone(c echo.Context) error {
	var (
		user models.User
		u    Register
		err  error
	)

	if err = c.Bind(user); err != nil {
		return err
	}

	db := orm.Conn
	if utility.IsValidAccount(user.Name) {
		err = db.Where("name = ?", user.Name).First(&u).Error
	}
	if err != nil {
		log.Logger.Error("User not found:", err)
	}

	if !utility.CompareHash([]byte(user.Pass), *u.Pass) {
		log.Logger.Error("Name and pass don't match:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidPass, "")

	}

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	sess.Set(general.SessionUserID, user.UserID)

	return c.JSON(errcode.ErrSucceed, nil)
}
