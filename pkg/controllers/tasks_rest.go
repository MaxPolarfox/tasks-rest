package controllers

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/MaxPolarfox/goTools/errors"
	"github.com/MaxPolarfox/goTools/mongoDB"
	"github.com/MaxPolarfox/tasks/pkg/types"
)

type RestTasksController struct {
	db DB
}

type DB struct {
	tasks mongoDB.Mongo
}

func NewRestTasksController(tasksCollection mongoDB.Mongo) *RestTasksController {
	return &RestTasksController{
		db: DB{
			tasks: tasksCollection,
		},
	}
}

// CreateTask  POST /rest/tasks
func (s *RestTasksController) CreateTask(rw http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	metricName := "RestTasksController.GetAllTasks"

	taskID := uuid.NewV4().String()

	createTaskBody := types.AddTaskReqBody{}
	err = json.NewDecoder(req.Body).Decode(&createTaskBody)
	if err != nil {
		log.Println(metricName+".Decode.Body", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	newTask := types.Task{
		ID:   taskID,
		Data: createTaskBody.Data,
	}

	// insert newTask to DB
	_, err = s.db.tasks.InsertOne(ctx, newTask)
	if err != nil {
		log.Println(metricName+"db.tasks.InsertOne", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	createdTaskRes := types.TaskIdResponse{ID: taskID}

	json, err := json.Marshal(createdTaskRes)
	if err != nil {
		log.Println(metricName+"json.Marshal", "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(json)
}

// GetAllTasks GET /rest/tasks
func (s *RestTasksController) GetAllTasks(rw http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	metricName := "RestTasksController.GetAllTasks"

	tasks := []types.Task{}

	filter := bson.M{}
	cursor, err := s.db.tasks.Find(ctx, filter)
	if err != nil {
		log.Println(metricName+".db.tasks.Find", "err", "err")
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		log.Println(metricName+".cursor.All", "err", "err")
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	json, err := json.Marshal(tasks)
	if err != nil {
		log.Println(metricName+".json.Marshal", "tasks", tasks, "err", err)
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(json)
}

// DeleteTask DELETE /rest/tasks/:id
func (s *RestTasksController) DeleteTask(rw http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	metricName := "RestTasksController.DeleteTask"

	params := httprouter.ParamsFromContext(ctx)

	taskID := params.ByName("id")

	filter := bson.M{"id": taskID}
	_, err = s.db.tasks.DeleteOne(ctx, filter)
	if err != nil {
		log.Println(metricName+".db.tasks.DeleteOne", "err", "err")
		errors.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}
