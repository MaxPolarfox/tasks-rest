package rest_client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gojektech/heimdall/v6"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gojektech/heimdall/v6/httpclient"

	goToolsClient "github.com/MaxPolarfox/goTools/client"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

type Client interface {
	AddTask(ctx context.Context, body types.AddTaskReqBody) (*types.TaskIdResponse, error)
	GetAllTasks(ctx context.Context) (*[]types.Task, error)
	DeleteTask(ctx context.Context, taskID string) error
}

type TasksClientImpl struct {
	URL    *url.URL
	client *httpclient.Client
}

func NewTasksClient(options goToolsClient.Options) Client {
	// First set a backoff mechanism. Constant backoff increases the backoff at a constant rate
	backoffInterval := 1000 * time.Millisecond
	// Define a maximum jitter interval. It must be more than 1*time.Millisecond
	maximumJitterInterval := 100 * time.Millisecond

	backoff := heimdall.NewConstantBackoff(backoffInterval, maximumJitterInterval)

	// Create a new retry mechanism with the backoff
	retrier := heimdall.NewRetrier(backoff)

	timeout := 10000 * time.Millisecond

	// Create a new client, sets the retry mechanism, and the number of times you would like to retry
	client := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(options.RetryCount),
	)

	parsedURL, err := url.Parse(options.URL)
	if err != nil {
		log.Fatalf("invalid URL", "err", err)
	}

	return &TasksClientImpl{
		URL:    parsedURL,
		client: client,
	}
}

func (c *TasksClientImpl) AddTask(ctx context.Context, body types.AddTaskReqBody) (*types.TaskIdResponse, error) {
	metricName := "TasksClientImpl.AddTask"

	createdTaskResponse := types.TaskIdResponse{}

	URL := goToolsClient.AppendToURL(c.URL, "/rest/tasks")
	urlString := URL.String()

	headers := http.Header{}

	js, err := json.Marshal(body)
	if err != nil {
		log.Println(metricName+"Marshal", "body", body, "err", err)
		return nil, err
	}

	res, err := c.client.Post(urlString, bytes.NewReader(js), headers)
	if err != nil {
		log.Println(metricName, "err", err, "u", urlString)
		return nil, err
	}
	defer res.Body.Close()

	if err := goToolsClient.ParseResponse(res, &createdTaskResponse); err != nil {
		return nil, err
	}
	return &createdTaskResponse, nil
}

func (c *TasksClientImpl) GetAllTasks(ctx context.Context) (*[]types.Task, error) {
	metricName := "TasksClientImpl.GetAllTasks"

	tasksResponse := []types.Task{}

	URL := goToolsClient.AppendToURL(c.URL, "/rest/tasks")
	urlString := URL.String()

	headers := http.Header{}

	res, err := c.client.Get(urlString, headers)
	if err != nil {
		log.Println(metricName, "err", err, "u", urlString)
		return nil, err
	}
	defer res.Body.Close()

	if err := goToolsClient.ParseResponse(res, &tasksResponse); err != nil {
		return nil, err
	}
	return &tasksResponse, nil
}

func (c *TasksClientImpl) DeleteTask(ctx context.Context, taskID string) error {
	metricName := "TasksClientImpl.AddTask"

	URL := goToolsClient.AppendToURL(c.URL, "/rest/tasks", taskID)
	urlString := URL.String()

	headers := http.Header{}

	res, err := c.client.Delete(urlString, headers)
	if err != nil {
		log.Println(metricName, "err", err, "u", urlString)
		return err
	}
	defer res.Body.Close()

	if err := goToolsClient.ParseResponse(res, nil); err != nil {
		return err
	}
	return nil
}
