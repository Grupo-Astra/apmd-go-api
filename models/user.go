// Package models contém as definições das estruturas de dados (ORM models).
package models

import "gorm.io/gorm"

// User representa a estrutura de um usuário no sistema.
//
// Ela é usada tanto para autenticação quanto para o armazenamento no banco de dados.
type User struct {
	gorm.Model

	// Username é o nome de usuário único utilizado para login.
	Username string `gorm:"unique;not null;index"`

	// Password armazena o hash bcrypt da senha do usuário.
	Password string `gorm:"not null"`
}
