package entity

import (
	"go-todolist/model"
	"go-todolist/utils/paginator"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryEntity interface {
	CreateCategory(category model.Category) (c model.Category, e error)
	// GetCategoryList(id int, name string) (categories []*model.Category)
	GetCategoryList(id int64, name string, page int64, limit int64) paginator.Page[model.Category]
	GetCategory(id int64) (res model.Category, err error)
	UpdateCategory(category model.Category) (c model.Category, e error)
	DeleteCategory(id int64) (c model.Category, e error)
}

type categoryConnection struct {
	connection *gorm.DB
}

func NewCategoryEntity(db *gorm.DB) CategoryEntity {
	return &categoryConnection{
		connection: db,
	}
}

func (db *categoryConnection) CreateCategory(category model.Category) (c model.Category, e error) {
	create := db.connection.Save(&category)
	if create.Error != nil {
		return category, create.Error
	}

	return category, nil
}

func (db *categoryConnection) GetCategoryList(id int64, name string, page int64, limit int64) paginator.Page[model.Category] {
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

func (db *categoryConnection) GetCategory(id int64) (category model.Category, err error) {
	res := db.connection.First(&category, "id=?", id)
	if res.Error == nil {
		return category, nil
	}

	return category, err
}

func (db *categoryConnection) UpdateCategory(category model.Category) (c model.Category, e error) {
	update := db.connection.Where("id = ?", category.ID).Updates(&category)
	if update.Error != nil {
		return category, update.Error
	}

	return category, nil
}

func (db *categoryConnection) DeleteCategory(id int64) (c model.Category, e error) {
	category := model.Category{}
	delete := db.connection.Delete(&category, id)
	if delete.Error != nil {
		return category, delete.Error
	}

	return category, nil
}
