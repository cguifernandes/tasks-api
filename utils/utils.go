package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// ParseValidationError recebe um erro do validator e retorna uma string com mensagens amigáveis.
func ParseValidationError(err error) string {
	if err == nil {
		return ""
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		msgs := []string{}
		for _, fieldErr := range validationErrors {
			msg := "Erro de validação no campo '" + fieldErr.Field() + "': "
			switch fieldErr.Tag() {
			case "required":
				msg += "obrigatório."
			case "max":
				msg += "tamanho máximo excedido."
			default:
				msg += "inválido."
			}
			msgs = append(msgs, msg)
		}
		return strings.Join(msgs, " ")
	}
	return err.Error()
}
