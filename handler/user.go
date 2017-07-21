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
 *	   Modify: 2017/07/19		  Ai Hao         添加用户登出
 *	   Modify: 2017/07/20         Zhang Zizhao   添加用户登录
 *    Modify: 2017/07/21          Xu Haosheng  更改用户信息
 */

package handler

import (
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"errors"
)

type Register struct {
	Mobile *string `json:"mobile" validate:"required,alphanum,min=6,max=30"`
	Pass   *string `json:"pass" validate:"required,alphanum,min=6,max=30"`
}

type GetPassword struct{
	Pass *string `json:"pass" validate:"required,alphanum,min=6,max=30"`
	NewPass *string `json:"newpass" validate:"required,alphanum,min=6,max=30"`
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

func LoginwithMobile(c echo.Context) error {
	var (
		user Register
		err  error
	)

	if err = c.Bind(&user); err != nil {
		log.Logger.Error("analysis creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	match := utility.IsValidAccount(*user.Mobile)
	if match == false {
		log.Logger.Error("err name format", err)

		return general.NewErrorWithMessage(errcode.ErrNameFormat, err.Error())
	}

	flag, userID, err := models.UserService.Login(user.Mobile, user.Pass)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("User not found:", err)

			return general.NewErrorWithMessage(errcode.ErrNamefound, err.Error())
		} else {
			log.Logger.Error("Mysql error:", err)

			return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
		}
	} else {
		if flag == false {
			log.Logger.Debug("Name and pass don't match:")

			return general.NewErrorWithMessage(errcode.ErrLoginRequired, errors.New("Name and pass don't match:").Error())
		}
	}

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
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

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetInfo(c echo.Context)  error {
	var (
		err    error
		Output models.UserInfo
	)

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	numberID := sess.Get(general.SessionUserID).(uint64)

	Output, err = models.UserService.GetInfo(numberID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("User information doesn't exist !", err)

			return general.NewErrorWithMessage(errcode.NoInformation, err.Error())
		}

		log.Logger.Error("Getting information exists errors", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Debug("have returned UserInformation.")

	return c.JSON(errcode.ErrSucceed, Output)
}

func ChangeMobilePassword (c echo.Context) error {
	var(
		password  	GetPassword
		userid    	uint64
		err       	error
		userpassword 	string
	)

	if err = c.Bind(&password); err != nil {
		log.Logger.Error("analysis creash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	sess := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	s := sess.Get(general.SessionUserID)
	userid = s.(uint64)

	userpassword, err = models.UserService.IsUserExist(userid)
	if err != nil {
		log.Logger.Error("User not found:", err)

		return general.NewErrorWithMessage(errcode.ErrNamefound, err.Error())
	}

	if !utility.CompareHash([]byte(userpassword), *password.Pass){
		log.Logger.Debug("Password doesn't match:", *password.Pass)

		return general.NewErrorWithMessage(errcode.ErrNamefound, errors.New("password").Error())
	}

	err = models.UserService.ChangeMobilePassword(password.NewPass, userid)
	if err != nil {
		log.Logger.Error("change faluse:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func ChangeUserinfo(c echo.Context) error {
	var (
		err  error
		info models.CUseInfo
	)

	if err = c.Bind(&info); err != nil {
		log.Logger.Error("Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	cuserID := session.Get(general.SessionUserID)
	id := cuserID.(uint64)
	log.Logger.Debug("id %d", id)

	err = models.UserService.ChangeUserInfo(info, id)
	if err != nil {
		log.Logger.Error("create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

/*func Changephone(c echo.Context) error {
=======
func Changephone(c echo.Context) error {
>>>>>>> 381348791139489194532cc0d0437d2b94e1b974
	var (
		err 	error
		m	models.Phone
	)
	if err = c.Bind(&m); err != nil {
		log.Logger.Error("ChangePhone crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}
	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	user := session.Get(general.SessionUserID).(uint64)

	err = models.UserService.ChangePhone(user,m.Phone)
	if err != nil {
		log.Logger.Error("changephone crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
<<<<<<< HEAD
}*/
