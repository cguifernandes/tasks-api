package routers

import (
	"example/tasks-api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Função responsável por gerar um JWT com as infos básicas do usuário
func generateJWT(userID, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    username,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// RegisterAuthRoutes organiza todas as rotas de autenticação (login e registro)
func RegisterAuthRoutes(router *gin.Engine, db *gorm.DB) {
	authGroup := router.Group("/auth")

	// Endpoint para registrar novos usuários
	authGroup.POST("/register", func(c *gin.Context) {
		var user models.User

		// Fazemos o bind para garantir que recebemos os dados certos
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos: " + err.Error(), "ok": false})
			return
		}

		// Checamos se já existe usuário com esse nome
		existing, err := models.GetByName(db, user.Name)
		if err == nil && existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Nome de usuário já cadastrado.", "ok": false})
			return
		}

		// Antes de salvar, nunca armazene senha em texto puro! Hash com bcrypt
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao gerar hash da senha.", "ok": false})
			return
		}
		user.Password = string(hash)

		// Tentativa de salvar de fato o usuário
		if err := user.SaveUser(db); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Erro ao criar o usuário: " + err.Error(), "ok": false})
			return
		}

		// Tudo certo! Usuário criado
		c.JSON(http.StatusCreated, gin.H{"user": user, "message": "Usuário registrado com sucesso!", "ok": true})
	})

	// Endpoint de login de usuário
	authGroup.POST("/login", func(c *gin.Context) {
		var login models.User

		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dados inválidos: " + err.Error(), "ok": false})
			return
		}

		// Busca o usuário pelo nome
		existing, err := models.GetByName(db, login.Name)
		if err != nil || existing == nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Usuário não encontrado.", "ok": false})
			return
		}

		// Verifica a senha: compara senha digitada com o hash
		if bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(login.Password)) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Senha incorreta.", "ok": false})
			return
		}

		// Gera token JWT válido por 72h
		token, err := generateJWT(existing.ID, existing.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Erro ao gerar token.", "ok": false})
			return
		}

		// Devolve o token para o frontend usar nas ações protegidas
		c.JSON(http.StatusOK, gin.H{"token": token, "user": existing, "message": "Login realizado com sucesso!", "ok": true})
	})
}
