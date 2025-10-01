package repository

import (
	"sport-events-backend/internal/config"
	"sport-events-backend/internal/models"
)

func CreateUser(user models.User) error {
	query := `
        INSERT INTO users (name, email, password, role) 
        VALUES ($1, $2, $3, $4) RETURNING id
    `
	_, err := config.DB.Exec(query, user.Name, user.Email, user.Password, user.Role)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := config.DB.Get(&user, query, email)
	return user, err
}

func GetUserByID(id int) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, role, created_at FROM users WHERE id = $1`
	err := config.DB.Get(&user, query, id)
	return user, err
}