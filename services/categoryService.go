package services

import (
	"go-todolist/entity"
	"go-todolist/model"
	"go-todolist/request"
	"go-todolist/utils/log"

	"github.com/mashingan/smapping"
)

type CategoryService interface {
	CreateCategory(category request.CategoryCreateOrUpdateRequest) (c model.Category, e error)
	UpdateCategory(category request.CategoryCreateOrUpdateRequest, id int64) (c model.Category, e error)
}

type categoryService struct {
	categoryEntity entity.CategoryEntity
}

func NewCategoryService(categoryEntity entity.CategoryEntity) CategoryService {
	return &categoryService{categoryEntity: categoryEntity}
}

func (s *categoryService) CreateCategory(category request.CategoryCreateOrUpdateRequest) (c model.Category, e error) {
	categoryToCreate := model.Category{}
	err := smapping.FillStruct(&categoryToCreate, smapping.MapFields(&category))
	if err != nil {
		log.Error("CreateCategory Failed map : " + err.Error())
		return categoryToCreate, err
	}

	res, resErr := s.categoryEntity.CreateCategory(categoryToCreate)
	if resErr != nil {
		return res, resErr
	}

	return res, nil
}

func (s *categoryService) UpdateCategory(category request.CategoryCreateOrUpdateRequest, id int64) (c model.Category, e error) {
	categoryToUpdate := model.Category{}
	err := smapping.FillStruct(&categoryToUpdate, smapping.MapFields(&category))
	if err != nil {
		log.Error("UpdateCategory Failed map : " + err.Error())
		return categoryToUpdate, err
	}

	categoryToUpdate.ID = id
	res, resErr := s.categoryEntity.UpdateCategory(categoryToUpdate)
	if resErr != nil {
		return res, resErr
	}

	return res, nil
}
