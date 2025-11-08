package models

import (
	"errors"
	"time"

	"example/tasks-api/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Task representa uma tarefa no sistema
type Task struct {
	ID          string    `json:"id" gorm:"primaryKey;type:varchar(36);unique;not null"`      // Identificador único, gerado por UUID
	Title       string    `json:"title" validate:"required,max=255" gorm:"size:255;not null"` // Título da tarefa (obrigatório)
	Description string    `json:"description" validate:"required,max=500" gorm:"size:500"`    // Descrição detalhada (obrigatória)
	Completed   bool      `json:"completed" gorm:"not null;default:false"`                    // Status de concluída
	UserID      string    `json:"user_id,omitempty" gorm:"type:varchar(36)"`                  // ID do usuário que criou a task
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`                    // Relação com o usuário
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`                           // Data/hora de criação
}

// Antes de salvar uma tarefa, gera o ID (se necessário) e valida os campos
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	// Validação de campos obrigatórios/limites via go-playground/validator
	validate := validator.New()
	err = validate.Struct(t)
	if err != nil {
		return errors.New(utils.ParseValidationError(err))
	}
	return nil
}

// Salva uma nova task no banco de dados
func (t *Task) SaveTask(db *gorm.DB) error {
	return db.Create(t).Error
}

// Atualiza uma task existente
func (t *Task) UpdateTask(db *gorm.DB) error {
	return db.Save(t).Error
}

// Remove uma task do banco
func (t *Task) DeleteTask(db *gorm.DB) error {
	return db.Delete(t).Error
}

// Recupera todas as tarefas cadastradas com informações do usuário
func GetAllTasks(db *gorm.DB) ([]Task, error) {
	var tasks []Task
	if err := db.Preload("User").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Busca uma tarefa por ID (retorna erro se não encontrar) com informações do usuário
func GetTaskById(db *gorm.DB, id string) (*Task, error) {
	var task Task
	if err := db.Preload("User").First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}
