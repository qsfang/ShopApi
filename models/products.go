package models

import (
	"time"
	"ShopApi/orm"
	"ShopApi/general"
)

type ProductServiceProvider struct {
}

var ProductService *ProductServiceProvider = &ProductServiceProvider{}

type Product struct {
	ID				uint64 `json:"id"`
	Name			string `json:"name"`
	Totalsale   	uint64 `json:"totalsale"`
	Categories		uint64 `json:"categories"`
	Price			float64 `json:"price"`
	Originalprice	float64 `json:"originalprice"`
	Status          uint64 `json:"status"`
	Size            string `json:"size"`
	Color           string `json:"color"`
	Imageid			uint64 `json:"imageid"`
	Imageids		string `json:"imageids"`
	Remark			string `json:"remark"`
	Detail			string `json:"detail"`
	Created			time.Time `json:"created"`
	Inventory		uint64 `json:"inventory"`
}

type CreatePro struct {
	Name			string `json:"name"`
	Categories		uint64 `json:"categories"`
	Price			float64 `json:"price"`
	Originalprice	float64 `json:"originalprice"`
	Size            string `json:"size"`
	Color           string `json:"color"`
	Imageid			uint64 `json:"imageid"`
	Imageids		string `json:"imageids"`
	Detail			string `json:"detail"`
	Inventory		uint64 `json:"inventory"`
}


func (Product) TableName() string {
	return "products"
}

func (ps *ProductServiceProvider) CreateP(pr CreatePro) error {
	pro := Product{
		Name:   		pr.Name,
		Categories:		pr.Categories,
		Price:			pr.Price,
		Originalprice:	pr.Originalprice,
		Status:         general.ProductOnsale,
		Size:           pr.Size,
		Color:          pr.Color,
		Imageid:		pr.Imageid,
		Imageids:		pr.Imageids,
		Detail:			pr.Detail,
		Created:		time.Now(),
		Inventory:		pr.Inventory,
	}

	db := orm.Conn

	err := db.Create(&pro).Error
	if err != nil {
		return err
	}

	return nil
}