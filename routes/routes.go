package routes

import (
	"appointment-scheduling/handlers"

	"github.com/gorilla/mux"
)

// InitializeRoutes initializes the API routes for the application
func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/appointments/available", handlers.GetAvailableAppointments).Methods("GET")
	router.HandleFunc("/appointments/create", handlers.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments/scheduled", handlers.ScheduledAppointments).Methods("GET")

	return router
}
