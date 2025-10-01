package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"sport-events-backend/internal/config"
	"sport-events-backend/internal/handlers"
	"sport-events-backend/internal/middleware"
)
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func main() {
	// Cargar variables de entorno
	err := godotenv.Load("D:/RA SPORT PROJECT/sport-events-backend/config/.env")
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	} else {
		log.Println("Archivo .env cargado correctamente")
	}
	// Iniciar conexión a la BD
	config.InitDB()
	jtw := os.Getenv("JWT_SECRET")
	if jtw == "" {
		log.Fatal("JWT_SECRET no está configurado")
	}
	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware) // protege todas las rutas /api/*

	// Solo organizers pueden crear eventos
	api.Handle("/events", middleware.RoleMiddleware("organizer")(http.HandlerFunc(handlers.CreateEvent))).Methods("POST")
	// Editar evento (organizer dueño)
	api.Handle("/events/{id}",middleware.RoleMiddleware("organizer")(http.HandlerFunc(handlers.UpdateEventHandler)),).Methods("PUT")
	// Eliminar evento (organizer dueño)
	api.Handle("/events/{id}",middleware.RoleMiddleware("organizer")(http.HandlerFunc(handlers.DeleteEventHandler)),).Methods("DELETE")
	// Solo organizers pueden ver inscritos
	api.Handle("/events/{id}/registrations", middleware.RoleMiddleware("organizer")(http.HandlerFunc(handlers.GetEventRegistrationsHandler))).Methods("GET")

	// Todos los autenticados pueden ver eventos
	api.HandleFunc("/events", handlers.GetEventsHandler).Methods("GET")

	// Solo runners pueden registrarse en eventos
	api.Handle("/events/{id}/register", middleware.RoleMiddleware("runner")(http.HandlerFunc(handlers.RegisterEventHandler))).Methods("POST")
	// Cancelar inscripción (solo runners)
	api.Handle("/events/{id}/register",	middleware.RoleMiddleware("runner")(http.HandlerFunc(handlers.CancelRegistrationHandler)),).Methods("DELETE")
	// Ver mis inscripciones (solo runners)
	api.Handle("/my-registrations",	middleware.RoleMiddleware("runner")(http.HandlerFunc(handlers.GetMyRegistrationsHandler)),).Methods("GET")


	// Cualquier usuario autenticado puede ver su propio perfil
	api.HandleFunc("/me", handlers.GetMeHandler).Methods("GET")


	// Rutas públicas
	router.HandleFunc("/auth/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/auth/login", handlers.LoginHandler).Methods("POST")


	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
