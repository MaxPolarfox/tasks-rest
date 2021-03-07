package tasks_rest

import (
	"github.com/MaxPolarfox/tasks-rest/pkg/controllers"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks-rest/pkg/types"
)

//TasksServiceGrpcImpl is a implementation of TasksService Grpc Service.
type Service struct {
	Router  *httprouter.Router
	options types.Options
	db      DB
}

type DB struct {
	tasks mongoDB.Mongo
}

//NewService returns the pointer to the Service.
func NewService(options types.Options, tasksCollection mongoDB.Mongo) *Service {

	tasksController := controllers.NewRestTasksController(tasksCollection)

	router := httprouter.New()

	// Routes
	router.HandlerFunc(http.MethodPost, "/rest/tasks/", tasksController.CreateTask)
	router.HandlerFunc(http.MethodGet, "/rest/tasks", tasksController.GetAllTasks)
	router.HandlerFunc(http.MethodDelete, "/rest/tasks/:id", tasksController.DeleteTask)

	return &Service{
		options: options,
		db:      DB{tasksCollection},
		Router:  router,
	}
}
