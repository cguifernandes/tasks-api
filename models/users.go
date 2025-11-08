package models

import (
	"errors"
	"time"

	"example/tasks-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User representa um usuário do sistema (login, senha, etc.)
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36);unique;not null"`            // UUID gerado automaticamente
	Name      string    `json:"name" validate:"required,max=255" gorm:"size:255;not null;unique"` // Nome de usuário (único e obrigatório)
	Password  string    `json:"-" validate:"required,max=255" gorm:"size:255;not null"`           // Senha (armazenada como hash, nunca exposta no JSON)
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`                                 // Data/hora do cadastro
}

// Antes de criar um usuário, gera UUID e valida os campos obrigatórios
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	validate := validator.New()
	err = validate.Struct(u)

	if err != nil {
		return errors.New(utils.ParseValidationError(err))
	}

	return nil
}

// Salvar o usuário no banco de dados
func (u *User) SaveUser(db *gorm.DB) error {
	return db.Create(u).Error
}

// Atualizar dados do usuário
func (u *User) UpdateUser(db *gorm.DB) error {
	return db.Save(u).Error
}

// Deletar usuário do banco
func (u *User) DeleteUser(db *gorm.DB) error {
	return db.Delete(u).Error
}

// Busca todos os usuários existentes na tabela
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Busca usuário pelo ID
func GetUserById(db *gorm.DB, id string) (*User, error) {
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Busca usuário pelo nome de usuário
func GetByName(db *gorm.DB, name string) (*User, error) {
	var user User
	if err := db.First(&user, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
