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
 *     Initial: 2017/07/18        Yusan Kurban
 *	   Modify: 2017/07/19         Sun Anxiang 添加用户登录
 *	   Modify: 2017/07/19		  Ai Hao 添加用户登出
 */

package handler

import (
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
	"github.com/astaxie/session"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Register struct {
	Mobile *string `json:"mobile" validate:"required,alphanum,min=6,max=30"`
	Pass   *string `json:"pass" validate:"required,alphanum,min=6,max=30"`
}

func Create(c echo.Context) error {
	var (
		err error
		u   Register
	)

	if err = c.Bind(&u); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.UserService.Create(u.Mobile, u.Pass)
	if err != nil {
		log.Logger.Error("create creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func Login(c echo.Context) error {
	var (
		err    error
		u      Register
		userID uint64
		sess   session.Session
	)

	if err = c.Bind(&u); err != nil {
		log.Logger.Error("Bind error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	flag, userID, err := models.UserService.Login(u.Mobile, u.Pass)
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			log.Logger.Error("User not found:", err)
		} else {

			log.Logger.Error("Mysql error:", err)
		}

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	} else {
		if flag == false {
			log.Logger.Error("Name and pass don't match:", err)

			return general.NewErrorWithMessage(errcode.ErrLoginRequired, err.Error())
		}
	}

	sess = utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	sess.Set(general.SessionUserID, userID)

	return c.JSON(errcode.ErrSucceed, nil)
}

func Logout(c echo.Context) error {
	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	err := sess.Delete(general.SessionUserID)

	if err != nil {
		log.Logger.Error("Logout with error", err)

		return general.NewErrorWithMessage(errcode.ErrDelete, err.Error())
	}

	log.Logger.Debug("i got here")
	return c.JSON(errcode.ErrSucceed, nil)
}

func GetInfo(c echo.Context) error {
	var (
		err error
		uu  models.UserInfo
		s   models.UserInfo
	)

	if err = c.Bind(&uu); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	s, err = models.GetInfo(uu.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("User information doesn't exist !", err)
			return general.NewErrorWithMessage(errcode.NoInformation, err.Error())
		} else {
			log.Logger.Error("Getting information exists errors", err)
			return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
		}
	}

	log.Logger.Debug("have returned UserInformation.")
	return c.JSON(errcode.ErrSucceed, s)
}
