package rest_client

import (
	"context"
	"net/url"

	"github.com/MaxPolarfox/tasks/pkg/types"
)

type Client interface {
	AddTask(ctx context.Context, data string) (*types.TaskIdResponse, error)
	GetAllTasks(ctx context.Context) (*[]types.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type TasksClientImpl struct {
	URL     *url.URL
}




