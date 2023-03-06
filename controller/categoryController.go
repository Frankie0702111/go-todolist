package controller

import (
	"go-todolist/entity"
	"go-todolist/request"
	"go-todolist/services"
	"go-todolist/utils/responses"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	Create(c *gin.Context)
	GetByList(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type categoryController struct {
	categoryService services.CategoryService
	categoryEntity  entity.CategoryEntity
}

func NewCategoryController(categoryService services.CategoryService, categoryEntity entity.CategoryEntity) CategoryController {
	return &categoryController{
		categoryService: categoryService,
		categoryEntity:  categoryEntity,
	}
}

func (h *categoryController) Create(c *gin.Context) {
	var input request.CategoryCreateOrUpdateRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	createCategory, createCategoryErr := h.categoryService.CreateCategory(input)
	if createCategoryErr != nil {
		match, _ := regexp.MatchString("Duplicate", createCategoryErr.Error())
		if match {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", "Error 1062: Duplicate entry "+input.Name, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		} else {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", createCategoryErr.Error(), nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
	}

	response := responses.SuccessResponse(http.StatusOK, "Create Success", createCategory)
	c.JSON(http.StatusOK, response)
	return
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
	response := responses.SuccessPageResponse(http.StatusOK, "Successfully get category list", category.CurrentPage, category.PageLimit, category.Total, category.Pages, category.Data)
	c.JSON(http.StatusOK, response)
	return
}

func (h *categoryController) Get(c *gin.Context) {
	var input request.CategoryGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := h.categoryEntity.GetCategory(input.Id)
	if categoryErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if category.ID == 0 {
		response := responses.SuccessResponse(http.StatusOK, "Record not found", nil)
		c.JSON(http.StatusOK, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Successfully get category", category)
	c.JSON(http.StatusOK, response)
	return
}

func (h *categoryController) Update(c *gin.Context) {
	var input request.CategoryCreateOrUpdateRequest
	var id request.CategoryGetRequest
	err := c.ShouldBindUri(&id)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := h.categoryEntity.GetCategory(id.Id)
	if categoryErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if category.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusNotFound, "Failed to process request", responses.RecordNotFound, nil)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	inputErr := c.ShouldBindJSON(&input)
	if inputErr != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", inputErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	updateCategory, updateCategoryErr := h.categoryService.UpdateCategory(input, id.Id)
	if updateCategoryErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", updateCategoryErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Update Success", updateCategory)
	c.JSON(http.StatusOK, response)
	return
}

func (h *categoryController) Delete(c *gin.Context) {
	var input request.CategoryGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	category, categoryErr := h.categoryEntity.GetCategory(input.Id)
	if categoryErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", categoryErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if category.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusNotFound, "Failed to process request", responses.RecordNotFound, nil)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	_, deleteCategoryErr := h.categoryEntity.DeleteCategory(input.Id)
	if deleteCategoryErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", deleteCategoryErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Delete Success", nil)
	c.JSON(http.StatusOK, response)
	return
}
