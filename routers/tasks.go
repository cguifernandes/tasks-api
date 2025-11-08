package routers

import (
	"example/tasks-api/models"
	"net/http"

	"example/tasks-api/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes gerencia todas as rotas de tarefa
func RegisterTasksRoutes(router *gin.Engine, db *gorm.DB) {
	// Cria um grupo para todas as rotas de /tasks
	tasksGroup := router.Group("/tasks")

	// Listar tarefas é público: qualquer usuário pode ver as tasks
	tasksGroup.GET("/", func(c *gin.Context) {
		tasks, err := models.GetAllTasks(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao buscar tarefas: " + err.Error(), "ok": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks, "message": "Tarefas listadas com sucesso", "ok": true})
	})

	// A partir daqui, tudo exige autenticação JWT

	// Criar uma nova task (apenas para usuários autenticados)
	tasksGroup.POST("/", middlewares.AuthMiddleware(), func(c *gin.Context) {
		var task models.Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos: " + err.Error(), "ok": false})
			return
		}

		// Pega o user_id do contexto (definido pelo middleware JWT)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Usuário não autenticado", "ok": false})
			return
		}
		task.UserID = userID.(string)

		if err := task.SaveTask(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao salvar a tarefa: " + err.Error(), "ok": false})
			return
		}
		// Task criada com sucesso
		c.JSON(http.StatusCreated, gin.H{"task": task, "message": "Tarefa criada com sucesso", "ok": true})
	})

	// Buscar detalhes de uma tarefa pelo ID (só autenticado consegue ver específica)
	tasksGroup.GET("/:id", middlewares.AuthMiddleware(), func(c *gin.Context) {
		id := c.Param("id")
		task, err := models.GetTaskById(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task não encontrada: " + err.Error(), "ok": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task": task, "message": "Tarefa encontrada", "ok": true})
	})

	// Atualizar uma task pelo ID (apenas autenticado)
	tasksGroup.PUT("/:id", middlewares.AuthMiddleware(), func(c *gin.Context) {
		id := c.Param("id")
		task, err := models.GetTaskById(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task não encontrada: " + err.Error(), "ok": false})
			return
		}
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos: " + err.Error(), "ok": false})
			return
		}
		if err := task.UpdateTask(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar a tarefa: " + err.Error(), "ok": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task": task, "message": "Tarefa atualizada com sucesso", "ok": true})
	})

	// Remover uma task pelo ID (só com token JWT)
	tasksGroup.DELETE("/:id", middlewares.AuthMiddleware(), func(c *gin.Context) {
		id := c.Param("id")
		task, err := models.GetTaskById(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task não encontrada: " + err.Error(), "ok": false})
			return
		}
		if err := task.DeleteTask(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao deletar a tarefa: " + err.Error(), "ok": false})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task": task, "message": "Tarefa deletada com sucesso", "ok": true})
	})
}
