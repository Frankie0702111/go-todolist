package entity

import (
	"go-todolist/model"
	"go-todolist/utils/paginator"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryEntity interface {
	// GetCategoryList(id int, name string) (categories []*model.Category)
	GetCategoryList(id int, name string, page int64, limit int64) paginator.Page[model.Category]
}

type categoryConnection struct {
	connection *gorm.DB
}

func NewCategoryEntity(db *gorm.DB) CategoryEntity {
	return &categoryConnection{
		connection: db,
	}
}

func (db *categoryConnection) GetCategoryList(id int, name string, page int64, limit int64) paginator.Page[model.Category] {
	// query := db.connection.Model(&categories).Preload(clause.Associations)
	// if len(name) > 0 {
	// 	query.Where("name = ?", name)
	// }
	// if id > 0 {
	// 	query.Where("id = ?", id)
	// }
	// query.Find(&categories)
	// return categories

	var categories []*model.Category
	query := db.connection.Model(&categories).Preload(clause.Associations)

	if id > 0 {
		query.Where("id = ?", id)
	}

	if len(name) > 0 {
		// query.Where(fmt.Sprintf("name like %q", (name + "%")))
		query.Where("name like ?", name+"%")
	}

	p := paginator.Page[model.Category]{CurrentPage: page, PageLimit: limit}
	p.SelectPages(query)

	return p
}
