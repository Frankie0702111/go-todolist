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

// @Summary "Create category"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."
// @Param	name			formData	string	true	"Category Name (maxLength: 100)"
// @Success 201 object responses.Response{errors=string,data=string} "Create Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/category [post]
func (h *categoryController) Create(c *gin.Context) {
	var input request.CategoryCreateOrUpdateRequest
	err := c.ShouldBind(&input)
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

	response := responses.SuccessResponse(http.StatusCreated, "Create Success", createCategory)
	c.JSON(http.StatusCreated, response)
	return
}

// @Summary "Category list"
// @Tags	"Category"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."
// @Param	id				query		integer	false	"Category ID"
// @Param	name			query		string	false	"Category Name (maxLength: 100)"
// @Param	page			query		integer	true	"Page (Please start from 1)"
// @Param	limit			query		integer	true	"Limit (Please start from 5 or 10)"
// @Success 200 object responses.PageResponse{errors=string,data=string} "Record not found || Successfully get category"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/category [get]
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

// @Summary	"Get a single category"
// @Tags	"Category"
// @Version	1.0
// @Produce	application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."
// @Param	id				path		integer	true	"Category ID"
// @Success	200 object responses.Response{errors=string,data=string} "Record not found || Successfully get category"
// @Failure	400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure	500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [get]
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

// @Summary	"Update a single category"
// @Tags	"Category"
// @Version	1.0
// @Produce	application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."
// @Param	id				path		integer	true	"Category ID"
// @Param	name			query		string	true	"Category Name (maxLength: 100)"
// @Success	200 object responses.Response{errors=string,data=string} "Update Success"
// @Failure	400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure	404 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure	500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [PATCH]
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

	inputErr := c.ShouldBind(&input)
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

// @Summary	"Delete a single category"
// @Tags	"Category"
// @Version	1.0
// @Produce	application/json
// @Param	Authorization	header		string	true	"example:Bearer token (Bearer+space+token)."
// @Param	id				path		integer	true	"Category ID"
// @Success	200 object responses.Response{errors=string,data=string} "Delete Success"
// @Failure	400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure	404 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure	500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/category/{id} [delete]
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
