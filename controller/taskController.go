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

	taskUuid := h.getUuid(input.UserID)
	createTask, createTaskErr := h.taskService.CreateTask(input, taskUuid)
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

	response := responses.SuccessResponse(http.StatusOK, "Create Success", createTask)
	c.JSON(http.StatusOK, response)
	return
}

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

func (h *taskController) Get(c *gin.Context) {
	var input request.TaskGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := h.taskEntity.GetTask(input.Id)
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	if task.ID == 0 {
		response := responses.SuccessResponse(http.StatusOK, "Record not found", nil)
		c.JSON(http.StatusOK, response)
		return
	}

	response := responses.SuccessResponse(http.StatusOK, "Successfully get category", task)
	c.JSON(http.StatusOK, response)
	return
}

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
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if task.ID == 0 {
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

func (h *taskController) Delete(c *gin.Context) {
	var input request.TaskGetRequest
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.IdInvalid, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	task, taskErr := h.taskEntity.GetTask(input.Id)
	if taskErr != nil {
		response := responses.ErrorsResponse(http.StatusInternalServerError, "Failed to process request", taskErr.Error(), nil)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if task.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusNotFound, "Failed to process request", responses.RecordNotFound, nil)
		c.AbortWithStatusJSON(http.StatusNotFound, response)
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
