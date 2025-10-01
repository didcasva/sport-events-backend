package models

import (
	"time"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // nunca se envía en JSON
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Registration struct {
	ID      int       `db:"id" gorm:"primaryKey" json:"id"`
	UserID  int       `db:"user_id" json:"user_id"`
	EventID int       `db:"event_id" json:"event_id"`
	Date    time.Time `db:"date" json:"date"`

	User User `json:"user"` // opcional para devolver info del usuario
}

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.StandardClaims
}
