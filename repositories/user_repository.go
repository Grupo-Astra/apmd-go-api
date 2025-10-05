// Package repositories contém as abstrações para acesso ao banco de dados.
package repositories

import (
	"github.com/Grupo-Astra/apmd-go-api/models"
	"gorm.io/gorm"
)

// UserRepositoryInterface define o contrato para as operações de
// usuários no banco de dados.
type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	ClearAll() error
}

// userRepository é a implementação concreta da UserRepositoryInterface.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository cria uma nova instância do repositório de usuário.
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

// Create cria um novo registro de usuário no banco de dados.
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByUsername busca um usuário pelo seu nome de usuário.
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// ClearAll remove todos os registros da tabela de usuários.
func (r *userRepository) ClearAll() error {
	return r.db.Exec("DELETE FROM users").Error
}
