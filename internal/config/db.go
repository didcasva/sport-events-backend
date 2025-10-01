package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// InitDB inicializa la conexión a PostgreSQL
func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Error al abrir conexión con la BD: %v", err)
	}

	// Configuración de pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Probar conexión
	if err := DB.Ping(); err != nil {
		log.Fatalf("❌ No se pudo conectar a la BD: %v", err)
	}

	log.Println("✅ Conectado a PostgreSQL")
}
