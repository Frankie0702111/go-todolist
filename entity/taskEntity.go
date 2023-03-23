package entity

import (
	"go-todolist/model"
	"go-todolist/utils/paginator"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskEntity interface {
	CreateTask(task model.Task) (c model.Task, e error)
	GetTaskList(id int64, user_id int64, title string, specify_datetime *time.Time, is_specify_time *bool, is_complete *bool, page int64, limit int64) paginator.Page[model.Task]
	GetTask(id int64) (task model.Task, err error)
	UpdateTask(task model.Task) (c model.Task, e error)
	DeleteTask(id int64) (c model.Task, e error)
	GetTaskImguuidByUserId(user_id int64) interface{}
}

type taskConnection struct {
	connection *gorm.DB
}

func NewTaskEntity(db *gorm.DB) TaskEntity {
	return &taskConnection{
		connection: db,
	}
}

func (db *taskConnection) CreateTask(task model.Task) (c model.Task, e error) {
	create := db.connection.Save(&task)
	if create.Error != nil {
		return task, create.Error
	}

	return task, nil
}

func (db *taskConnection) GetTaskList(id int64, user_id int64, title string, specify_datetime *time.Time, is_specify_time *bool, is_complete *bool, page int64, limit int64) paginator.Page[model.Task] {
	var tasks []*model.Task
	query := db.connection.Model(&tasks).Preload(clause.Associations)

	if id > 0 {
		query.Where("id = ?", id)
	}

	if user_id > 0 {
		query.Where("user_id = ?", user_id)
	}

	if len(title) > 0 {
		query.Where("title like ?", title+"%")
	}

	if specify_datetime != nil {
		query.Where("specify_datetime = ?", specify_datetime)
	}

	if is_specify_time != nil {
		query.Where("is_specify_time = ?", is_specify_time)
	}

	if is_complete != nil {
		query.Where("is_complete = ?", is_complete)
	}

	p := paginator.Page[model.Task]{CurrentPage: page, PageLimit: limit}
	p.SelectPages(query)

	return p
}

func (db *taskConnection) GetTask(id int64) (task model.Task, err error) {
	res := db.connection.Preload("Category").First(&task, "id=?", id)
	if res.Error == nil {
		return task, nil
	}

	return task, err
}

func (db *taskConnection) UpdateTask(task model.Task) (c model.Task, e error) {
	update := db.connection.Where("id=?", task.ID).Updates(&task)
	if update.Error != nil {
		return task, update.Error
	}

	return task, nil
}

func (db *taskConnection) DeleteTask(id int64) (c model.Task, e error) {
	task := model.Task{}
	delete := db.connection.Delete(&task, id)
	if delete.Error != nil {
		return task, delete.Error
	}

	return task, nil
}

func (db *taskConnection) GetTaskImguuidByUserId(user_id int64) interface{} {
	task := model.Task{}
	res := db.connection.Where("img_uuid IS NOT NULL AND user_id = ?", user_id).First(&task)
	if res.Error != nil {
		return nil
	}

	return task.ImgUuid
}
