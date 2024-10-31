package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"go-backend/internal/db"
	"go-backend/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type Container struct {
	q *db.Queries
}

func NewContainer(database *sql.DB) (*Container, error) {
	return &Container{q: db.New(database)}, nil
}

// UsersPost - Create a new user
func (c *Container) UsersPost(ctx echo.Context) error {
	var newUser models.User
	if err := ctx.Bind(&newUser); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	u, err := c.q.CreateUser(ctx.Request().Context(), db.CreateUserParams{
		Uid:   newUser.UID,
		Name:  newUser.Name,
		Email: newUser.Email,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]int{"id": int(u.ID)})
}

// UsersGet - Get a list of users
func (c *Container) UsersGet(ctx echo.Context) error {

	users, err := c.q.GetUsers(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	send_users := []models.User{}
	for _, user := range users {
		send_users = append(send_users, models.User{
			ID:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return ctx.JSON(http.StatusOK, send_users)
}

// TasksGet - Get a list of tasks
func (c *Container) TasksGet(ctx echo.Context) error {
	tasks, err := c.q.GetTasks(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	send_tasks := []models.Task{}
	for _, task := range tasks {
		send_tasks = append(send_tasks, models.Task{
			ID:             int(task.ID),
			Title:          task.Title,
			TaskTime:       task.TaskTime.String(),
			TaskFinishDate: task.TaskFinishDate.String(),
			UserID:         int(task.UserID),
		})
	}
	return ctx.JSON(http.StatusOK, send_tasks)
}

// TasksPost - Create a new task
func (c *Container) TasksPost(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		fmt.Println("No user found in context")
		return errors.New("User not found in context")
	}

	user_data, ok := u.(models.User)
	if !ok {
		fmt.Printf("Failed to assert type: %T\n", u)
		return errors.New("Failed to assert user data")
	}
	var newTask models.NewTask
	if err := ctx.Bind(&newTask); err != nil {
		log.Printf("Failed to bind request body to NewTask: %v", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task_time, err := time.Parse("15:04:05", newTask.TaskTime)
	if err != nil {
		log.Printf("Failed to bind request body to NewTask: %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	task_finish_date, err := time.Parse("2006-01-02T15:04:05.000Z", newTask.TaskFinishDate)
	if err != nil {
		log.Printf("Failed to bind request body to NewTask: %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	task, err := c.q.CreateTask(ctx.Request().Context(), db.CreateTaskParams{
		Title:          newTask.Title,
		TaskTime:       task_time,
		TaskFinishDate: task_finish_date,
		UserID:         int32(user_data.ID),
	})
	if err != nil {
		log.Printf("Failed to bind request body to NewTask: %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	db_tags, err := c.q.GetTags(ctx.Request().Context())
	if err != nil {
		log.Printf("Failed to bind request body to NewTask: %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, tag := range newTask.Tags {
		if isContainTag(tag, db_tags) {
			tag_id, err := c.q.GetTagIdByName(ctx.Request().Context(), tag)
			if err != nil {
				log.Printf("Failed to bind request body to NewTask: %v", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
			_, err = c.q.CreateTaskTag(ctx.Request().Context(), db.CreateTaskTagParams{
				TaskID: int32(task.ID),
				TagID:  tag_id,
			})
			if err != nil {
				log.Printf("Failed to bind request body to NewTask: %v", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
		} else {
			new_tag, err := c.q.CreateTag(ctx.Request().Context(), tag)
			if err != nil {
				log.Printf("Failed to bind request body to NewTask: %v", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
			_, err = c.q.CreateTaskTag(ctx.Request().Context(), db.CreateTaskTagParams{
				TaskID: int32(task.ID),
				TagID:  int32(new_tag.ID),
			})
			if err != nil {
				log.Printf("Failed to bind request body to NewTask: %v", err)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}

	return ctx.JSON(http.StatusCreated, map[string]int{"id": int(task.ID)})
}

func isContainTag(tag string, db_tags []db.Tag) bool {
	for _, db_tag := range db_tags {
		if db_tag.Name == tag {
			return true
		}
	}
	return false
}

// TasksTaskIdDelete - Delete a task
func (c *Container) TasksTaskIdDelete(ctx echo.Context) error {
	id := ctx.Param("taskId")
	taskID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	_, err = c.q.DeleteTask(ctx.Request().Context(), int32(taskID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusNoContent)
}

// TasksTaskIdGet - Get a task by ID
func (c *Container) TasksTaskIdGet(ctx echo.Context) error {
	id := ctx.Param("taskId")

	taskID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	task, err := c.q.GetTaskByID(ctx.Request().Context(), int32(taskID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, int(task.ID))
}

// TasksTaskIdPut - Update a task
func (c *Container) TasksTaskIdPut(ctx echo.Context) error {
	id := ctx.Param("taskId")
	var updateTask models.UpdateTask
	if err := ctx.Bind(&updateTask); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	taskID, err := strconv.ParseInt(id, 10, 32)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := c.q.UpdateTask(ctx.Request().Context(), db.UpdateTaskParams{
		ID:             int32(taskID),
		Title:          updateTask.Title,
		TaskTime:       time.Time{},
		TaskFinishDate: time.Time{},
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	send_task := models.Task{
		ID:             int(task.ID),
		Title:          task.Title,
		TaskTime:       task.TaskTime.String(),
		TaskFinishDate: task.TaskFinishDate.String(),
		UserID:         int(task.UserID),
	}
	return ctx.JSON(http.StatusOK, send_task)
}

// get a group of tasks by tag
func (c *Container) TasksTagGet(ctx echo.Context) error {
	u := ctx.Get("user")
	if u == nil {
		fmt.Println("No user found in context")
		return errors.New("User not found in context")
	}

	user_data, ok := u.(models.User)
	if !ok {
		fmt.Printf("Failed to assert type: %T\n", u)
		return errors.New("Failed to assert user data")
	}

	query_start_date := ctx.QueryParam("start_date")
	query_end_date := ctx.QueryParam("end_date")

	start_time, err := time.Parse("2006-01-02", query_start_date)
	if err != nil {
		start_time, err = time.Parse("2006-01-02", "0001-01-01")
	}
	end_time, err := time.Parse("2006-01-02", query_end_date)

	if err != nil {
		end_time, err = time.Parse("2006-01-02", "9999-12-31")
	}

	task_tags, err := c.q.GetTaskTags(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// make a map task list by tag id
	tag_id_list := []int{}
	tag_task_dict := make(map[int][]int)
	for _, task_tag := range task_tags {
		tag_id_list = append(tag_id_list, int(task_tag.TagID))
		tag_task_dict[int(task_tag.TagID)] = append(tag_task_dict[int(task_tag.TagID)], int(task_tag.TaskID))
	}
	// unique tag id list
	tag_id_list = unique(tag_id_list)

	// sum task_time by tag
	sum_task_time_dict := make(map[string]int)
	for _, tag_id := range tag_id_list {
		task_ids := tag_task_dict[tag_id]
		sum_task_time := 0
		tag_name, err := c.q.GetTagNameByID(ctx.Request().Context(), int32(tag_id))
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		for _, task_id := range task_ids {
			task, err := c.q.GetTaskByID(ctx.Request().Context(), int32(task_id))
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
			if isTimeCalcTaeget(task.TaskFinishDate, start_time, end_time, user_data.ID, task.UserID) {
				sum_task_time += int(task.TaskTime.Hour())*60 + int(task.TaskTime.Minute())
			}
		}
		sum_task_time_dict[tag_name] = sum_task_time

	}

	return ctx.JSON(http.StatusOK, sum_task_time_dict)
}

func (c *Container) getOrCreateUser(ctx echo.Context, claims jwt.MapClaims) (models.User, error) {
	uid, ok := claims["user_id"].(string)
	if !ok {
		return models.User{}, errors.New("Failed to assert user_id")
	}
	user, err := c.q.GetUserByUID(ctx.Request().Context(), uid)
	if err != nil {
		// new user
		newUser := &models.User{
			UID:   uid,
			Name:  "",
			Email: "",
		}

		//
		if name, ok := claims["name"].(string); ok {
			newUser.Name = name
		}
		if email, ok := claims["email"].(string); ok {
			newUser.Email = email
		}

		//
		user, err = c.q.CreateUser(ctx.Request().Context(), db.CreateUserParams{
			Uid:   newUser.UID,
			Name:  newUser.Name,
			Email: newUser.Email,
		})
		if err != nil {
			return models.User{}, err
		}
	}

	return models.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		UID:   user.Uid,
	}, err

}

func isTimeCalcTaeget(time time.Time, start_time time.Time, end_time time.Time, user_id int, i int32) bool {

	if time.After(start_time) && time.Before(end_time) && user_id == int(i) {
		return true
	}
	return false
}

func unique(tag_id_list []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range tag_id_list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
