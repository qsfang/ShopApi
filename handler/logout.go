package handler

import (
	"ShopApi/general/errcode"
	"ShopApi/utility"
	"net/http"

	"github.com/labstack/echo"
	_"github.com/go-sql-driver/mysql"

	"ShopApi/general"
)

func Logout(c echo.Context) error {
	status := errcode.ErrSucceed

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	err := sess.Delete("login")

	if err != nil {
		return general.NewErrorWithMessage(errcode.ErrDelete, err.Error())
	}
	status = http.StatusOK

	return c.JSONPretty(http.StatusOK, status, " ")
}
