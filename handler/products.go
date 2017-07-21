package handler

import (
	"ShopApi/log"
	"ShopApi/models"
	"github.com/labstack/echo"
	"ShopApi/general"
	"ShopApi/general/errcode"
)

//名称name，totalsale  ，类型categories，价格price，原价originalprice，
// 状态status，尺码siez，颜色color,封面图片imageid，图片集imageids，评论remark,
//详细信息 detail ，创建日期 created，存货量inventory

func CreateP(c echo.Context) error {
	var (
		err 	error
		p		models.CreatePro
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

func ChangeProStatus(c echo.Context) error {
	var(
		err		error
		pro		models.ChangePro
	)

	if err = c.Bind(&pro); err != nil {
		log.Logger.Error("Change crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	err = models.ProductsService.ChangeProStatus(pro)
	if err != nil {
		log.Logger.Error("change crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}
