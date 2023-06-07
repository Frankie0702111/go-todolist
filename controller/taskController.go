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

type TaskController interface {
	Create(c *gin.Context)
	GetByList(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type taskController struct {
	taskService services.TaskService
	taskEntity  entity.TaskEntity
}

func NewTaskController(taskService services.TaskService, taskEntity entity.TaskEntity) TaskController {
	return &taskController{
		taskService: taskService,
		taskEntity:  taskEntity,
	}
}

func (h *taskController) getUuid(user_id int64) interface{} {
	task := h.taskEntity.GetTaskImguuidByUserId(user_id)
	return task
}

// @Summary "Create task"
// @Tags	"Task"
// @Version 1.0
// @Accept	multipart/form-data
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	user_id				formData	integer	true	"User ID"											minimum(1)
// @Param	category_id			formData	integer	true	"Category ID"										minimum(1)
// @Param	title				formData	string	true	"Title"												maxLength(100)
// @Param	note				formData	string	false	"Note"
// @Param	url					formData	string	false	"Url"
// @Param	image				formData	file	false	"Image"
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	priority			formData	integer	true	"Priority"											Enums(1, 2, 3) default(1)
// @Param	is_complete			formData	boolean	false	"Is Complete"										default(false)
// @Success 201 object responses.Response{errors=string,data=string} "Create Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/task [post]
func (h *taskController) Create(c *gin.Context) {
	var input request.TaskCreateRequest
	err := c.ShouldBind(&input)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if input.Image != nil {
		if len(input.Image.Filename) > 100 {
			response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.ImageFileNameLimitOf100, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		if input.Image.Size > (5 << 20) {
			response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.ImageFileSizeLimitOf5MB, nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}

	createTask, createTaskErr := h.taskService.CreateTask(input, nil)
	if createTaskErr != nil {
		match, _ := regexp.MatchString("Duplicate", createTaskErr.Error())
		if match {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", "Error 1062: Duplicate entry "+input.Title, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		} else {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", createTaskErr.Error(), nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
	}

	response := responses.SuccessResponse(http.StatusCreated, "Create Success", createTask)
	c.JSON(http.StatusCreated, response)
	return
}

// @Summary "Task list"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	id					formData	integer	false	"Task ID"											minimum(1)
// @Param	user_id				formData	integer	false	"User ID"											minimum(1)
// @Param	title				formData	string	false	"Title"												maxLength(100)
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	is_complete			formData	boolean	false	"Is Complete"
// @Param	page				query		integer	true	"Page"												minimum(1) default(1)
// @Param	limit				query		integer	true	"Limit"												minimum(2) default(5)
// @Success 200 object responses.PageResponse{errors=string,data=string} "Successfully get task list"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/task [get]
func (h *taskController) GetByList(c *gin.Context) {
	var input request.TaskGetListRequest
	err := c.ShouldBind(&input)
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task := h.taskEntity.GetTaskList(input.Id, input.UserID, input.Title, input.SpecifyDatetime, input.IsSpecifyTime, input.IsComplete, input.Page, input.Limit)
	response := responses.SuccessPageResponse(http.StatusOK, "Successfully get task list", task.CurrentPage, task.PageLimit, task.Total, task.Pages, task.Data)
	c.JSON(http.StatusOK, response)
	return
}

// @Summary "Get a single task"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Task ID"										minimum(1)
// @Success 200 object responses.Response{errors=string,data=string} "Record not found || Successfully get task"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [get]
func (h *taskController) Get(c *gin.Context) {
	var input request.TaskGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := h.taskEntity.GetTask(input.Id)
	if task.ID == 0 {
		response := responses.SuccessResponse(http.StatusOK, "Record not found", nil)
		c.JSON(http.StatusOK, response)
		return
	}
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Successfully get task", task)
	c.JSON(http.StatusOK, response)
	return
}

// @Summary "Update a single task"
// @Tags	"Task"
// @Version 1.0
// @Accept	multipart/form-data
// @Produce application/json
// @Param	Authorization		header		string	true	"example:Bearer token (Bearer+space+token)."		default(Bearer )
// @Param	id					path		integer	true	"Task ID"											minimum(1)
// @Param	category_id			formData	integer	false	"Category ID"										minimum(1)
// @Param	title				formData	string	false	"Title"												maxLength(100)
// @Param	note				formData	string	false	"Note"
// @Param	url					formData	string	false	"Url"
// @Param	image				formData	file	false	"Image"
// @Param	specify_datetime	formData	string	false	"Specify Datetime (DateTime: 2006-01-02 15:04:05)"
// @Param	is_specify_time		formData	boolean	false	"Is Specify Time"
// @Param	priority			formData	integer	true	"Priority"											Enums(1, 2, 3)
// @Param	is_complete			formData	boolean	false	"Is Complete"
// @Success 200 object responses.Response{errors=string,data=string} "Update Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [PATCH]
func (h *taskController) Update(c *gin.Context) {
	var input request.TaskUpdateRequest
	var id request.TaskGetRequest
	err := c.ShouldBindUri(&id)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := h.taskEntity.GetTask(id.Id)
	if task.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusNotFound, "Failed to process request", responses.RecordNotFound, nil)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	inputErr := c.ShouldBind(&input)
	if inputErr != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", inputErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	taskUuid := h.getUuid(task.UserID)
	updateTask, updateTaskErr := h.taskService.UpdateTask(input, id.Id, task.UserID, task.Img, taskUuid)
	if updateTaskErr != nil {
		match, _ := regexp.MatchString("Duplicate", updateTaskErr.Error())
		if match {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", "Error 1062: Duplicate entry "+input.Title, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		} else {
			response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", updateTaskErr.Error(), nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
	}

	response := responses.SuccessResponse(http.StatusOK, "Update Success", updateTask)
	c.JSON(http.StatusOK, response)
	return
}

// @Summary "Delete a single task"
// @Tags	"Task"
// @Version 1.0
// @Produce application/json
// @Param	Authorization	header	string	true	"example:Bearer token (Bearer+space+token)."	default(Bearer )
// @Param	id				path	integer	true	"Task ID"										minimum(1)
// @Success 200 object responses.Response{errors=string,data=string} "Delete Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 404 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/task/{id} [delete]
func (h *taskController) Delete(c *gin.Context) {
	var input request.TaskGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := h.taskEntity.GetTask(input.Id)
	if task.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusNotFound, "Failed to process request", responses.RecordNotFound, nil)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	_, deleteTaskErr := h.taskEntity.DeleteTask(input.Id)
	if deleteTaskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", deleteTaskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Delete Success", nil)
	c.JSON(http.StatusOK, response)
	return
}
