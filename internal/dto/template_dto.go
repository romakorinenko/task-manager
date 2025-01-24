package dto

import (
	"github.com/romakorinenko/task-manager/internal/repository"
)

type TasksWithLoginTemplateData struct {
	Tasks []repository.TaskWithLogin
}

type UsersTemplateData struct {
	Users []repository.User
}
