package controller

import (
	"go-todolist/entity"
	"go-todolist/request"
	"go-todolist/utils/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	GetByList(c *gin.Context)
}

type categoryController struct {
	categoryEntity entity.CategoryEntity
}

func NewCategoryController(categoryEntity entity.CategoryEntity) CategoryController {
	return &categoryController{
		categoryEntity: categoryEntity,
	}
}

func (h *categoryController) GetByList(c *gin.Context) {
	var input request.CategoryGetListRequest
	err := c.ShouldBind(&input)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	category := h.categoryEntity.GetCategoryList(input.Id, input.Name, input.Page, input.Limit)
	response := responses.SuccessResponse(http.StatusOK, "Successfully get category list", category)
	c.JSON(http.StatusOK, response)
	return
}
