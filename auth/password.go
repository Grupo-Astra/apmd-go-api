// Package auth contém a lógica de negócio para autenticação,
// como manipulação de senhas e tokens JWT.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword gera um hash bcrypt a partir de uma senha.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compara uma senha em texto plano com um hash bcrypt existente.
//
// Retorna true se a senha corresponder ao hash, e false caso contrário.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}
