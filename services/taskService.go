package services

import (
	"fmt"
	"go-todolist/entity"
	"go-todolist/model"
	"go-todolist/request"
	"go-todolist/utils/log"

	"github.com/gofrs/uuid"
	"github.com/mashingan/smapping"
)

type TaskService interface {
	CreateTask(task request.TaskCreateRequest, img_uuid interface{}) (c model.Task, e error)
	UpdateTask(task request.TaskUpdateRequest, id int64, user_id int64, img string, img_uuid interface{}) (c model.Task, e error)
}

type taskService struct {
	taskEntity entity.TaskEntity
	s3Entity   entity.S3Entity
}

func NewTaskService(taskEntity entity.TaskEntity, s3Entity entity.S3Entity) TaskService {
	return &taskService{
		taskEntity: taskEntity,
		s3Entity:   s3Entity,
	}
}

func (s *taskService) CreateTask(task request.TaskCreateRequest, img_uuid interface{}) (c model.Task, e error) {
	taskToCreate := model.Task{}
	err := smapping.FillStruct(&taskToCreate, smapping.MapFields(&task))
	if err != nil {
		log.Error("CreateTask Failed map : " + err.Error())
		return taskToCreate, err
	}

	if task.Image != nil {
		uuidV4 := ""
		if img_uuid == nil {
			uuidV4Ojb, uuidV4Err := uuid.NewV4()
			if uuidV4Err != nil {
				return taskToCreate, uuidV4Err
			}
			uuidV4 = uuidV4Ojb.String()
		} else {
			uuidV4 = fmt.Sprintf("%s", img_uuid)
		}

		s3Res, s3ResErr := s.s3Entity.FileUpload(task.Image, uuidV4)
		if s3ResErr != nil {
			return taskToCreate, s3ResErr
		}
		taskToCreate.Img = task.Image.Filename
		taskToCreate.ImgLink = s3Res.Location
		taskToCreate.ImgUuid = uuidV4
	}

	res, resErr := s.taskEntity.CreateTask(taskToCreate)
	if resErr != nil {
		return res, resErr
	}

	return res, nil
}

func (s *taskService) UpdateTask(task request.TaskUpdateRequest, id int64, user_id int64, img string, img_uuid interface{}) (c model.Task, e error) {
	taskToUpdate := model.Task{}
	err := smapping.FillStruct(&taskToUpdate, smapping.MapFields(&task))
	if err != nil {
		log.Error("CreateTask Failed map : " + err.Error())
		return taskToUpdate, err
	}

	if task.Image != nil {
		uuidV4 := ""
		if img_uuid == nil {
			uuidV4Ojb, uuidV4Err := uuid.NewV4()
			if uuidV4Err != nil {
				return taskToUpdate, uuidV4Err
			}
			uuidV4 = uuidV4Ojb.String()
		} else {
			uuidV4 = fmt.Sprintf("%s", img_uuid)
		}

		if task.Image.Filename != img {
			s3RemoveErr := s.s3Entity.FileRemove(img, uuidV4)
			if s3RemoveErr != nil {
				return taskToUpdate, s3RemoveErr
			}
		}
		s3UploadRes, s3UploadResErr := s.s3Entity.FileUpload(task.Image, uuidV4)
		if s3UploadResErr != nil {
			return taskToUpdate, s3UploadResErr
		}
		taskToUpdate.Img = task.Image.Filename
		taskToUpdate.ImgLink = s3UploadRes.Location
		taskToUpdate.ImgUuid = uuidV4
	}

	taskToUpdate.ID = id
	taskToUpdate.UserID = user_id
	res, resErr := s.taskEntity.UpdateTask(taskToUpdate)
	if resErr != nil {
		return res, resErr
	}

	return res, nil
}
