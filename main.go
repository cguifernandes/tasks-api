package main

import (
	"log"

	"example/tasks-api/models"
	"example/tasks-api/routers"

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

	// Aponta rotas de tasks e autenticação
	routers.RegisterRoutes(router, db)
	routers.RegisterAuthRoutes(router, db)

	// Sobe a API na porta 8080
	router.Run(":8080")
}
