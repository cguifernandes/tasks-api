package main

import (
	"log"

	"example/tasks-api/models"
	"example/tasks-api/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Inicializa banco de dados e faz migrate automática das structs
func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Garante que as tabelas estejam alinhadas com os models sempre que subir a app
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.User{})
	return db, nil
}

func main() {
	// Conecta e prepara banco
	db, err := initDB()
	if err != nil {
		log.Fatalf("Erro ao inicializar o banco de dados: %v", err)
		return
	}

	// Cria o router principal do Gin
	router := gin.Default()

	// Configura CORS para permitir requisições de outras origens
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:5500",
			"http://localhost:8000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Aponta rotas de tasks e autenticação
	routers.RegisterTasksRoutes(router, db)
	routers.RegisterAuthRoutes(router, db)

	// Sobe a API na porta 8080
	router.Run(":8080")
}
