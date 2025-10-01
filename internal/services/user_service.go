package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"sport-events-backend/internal/models"
	"sport-events-backend/internal/repository"
)

func RegisterUser(name, email, password, role string) error {
	// encriptar password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     role,
	}

	return repository.CreateUser(user)
}

func AuthenticateUser(email, password string) (models.User, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return models.User{}, errors.New("usuario no encontrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("contraseña inválida")
	}

	return user, nil
}
