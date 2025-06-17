package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// структура задачи 
type Task struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"` // тут мы обозначаем как будит писаться это поле в json
	Description string    `json:"description`
	IsComplete  bool      `json:"is_complete"`
	CreatedAt   time.Time `json:""created_at`
}

var tasks = make(map[uuid.UUID]Task) // типа хранилище

func main() {

	r := gin.Default() // так называемый ролтер ( отправляет запросы куда надо ) энд поинт

	// создать новую задачу
	r.POST("/tasks/", HandleCreateTask) //  создание задачи
	// получить все задачи
	r.GET("/tasks/", HandleGetAllTasks) // получить все задачи
	// поменять значения выполненасти
	r.POST("/tasks/:taskId/", HandleToggleComplete) // переключить статус
	r.GET("/tasks/:taskId", HandleGetTaskByID)   // получить задачу 
	r.DELETE("/tasks/:taskId", HandleDeleteTask) //  удалить задачу
	r.PUT("/tasks/:taskId", HandleUpdateTask) // изменить задачу

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}

func HandleCreateTask(c *gin.Context) {
	var task Task
	err := json.NewDecoder(c.Request.Body).Decode(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, "у тебя не правельный json"+err.Error())
	}

	id := uuid.New()
	task.Id = id
	task.CreatedAt = time.Now()
	tasks[id] = task
	c.JSON(http.StatusOK, nil)
}

func HandleGetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func HandleToggleComplete(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ты указал не правельные id задачи ")
	}

	task := tasks[id]
	task.IsComplete = !task.IsComplete
	tasks[id] = task

	c.JSON(http.StatusOK, task)
}

func HandleGetTaskByID(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Неправильный ID задачи")
		return
	}

	task := tasks[id]

	if task.Id == uuid.Nil {
		c.JSON(http.StatusNotFound, "Задача не найдена")
		return
	}

	c.JSON(http.StatusOK, task)
}

func HandleDeleteTask(c *gin.Context) {
	idStr := c.Param("taskId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Неправильный ID задачи")
		return
	}

	task := tasks[id]
	delete(tasks, id)

	if task.Id == uuid.Nil {
		c.JSON(http.StatusOK, "Задачи и так не существовало")
	} else {
		c.JSON(http.StatusOK, task)
	}
}

///  я не смог его доделать ( голова взорволась)
func HandleUpdateTask(c *gin.Context) {
    // получаем ID задачи из URL
    taskId := c.Param("taskId")

    // ищем задачу по ID
    task, err := getTaskById(taskId)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Некорректный ID задачи"}) // 404 более уместен
        return
    }

    // 3. Читаем новые данные из тела запроса
    var updates struct {
        Title       string `json:"title"`       
        Description string `json:"description"`
        IsComplete  bool   `json:"is_complete"`
    }