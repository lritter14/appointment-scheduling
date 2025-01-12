package handlers

import (
	"appointment-scheduling/models"
	"appointment-scheduling/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

/*
GetAvailableAppointments GET request to get a list of available appointment times for a trainer between two dates.

Path: /appointments/available
Parameters:
  - trainer_id (int): ID of the trainer for the appointment
  - starts_at (string): Start date in format "2019-01-25"
  - ends_at (string): End date in format "2019-01-25"

Returns:
  - 200 OK with list of available appointment start times
    [
    "2019-01-25T09:00:00-08:00",
    "2019-01-25T09:30:00-08:00",
    "2019-01-25T10:00:00-08:00",
    "2019-01-25T10:30:00-08:00",
    "2019-01-25T11:00:00-08:00",
    ]
  - 400 Bad Request if request body is invalid
*/
func GetAvailableAppointments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	trainerID := query.Get("trainer_id")
	startsAt := query.Get("starts_at")
	endsAt := query.Get("ends_at")

	if trainerID == "" {
		http.Error(w, "Missing trainer_id", http.StatusBadRequest)
		return
	}
	if startsAt == "" {
		http.Error(w, "Missing starts_at", http.StatusBadRequest)
		return
	}
	if endsAt == "" {
		http.Error(w, "Missing ends_at", http.StatusBadRequest)
		return
	}

	trainerIDInt, err := strconv.Atoi(trainerID)
	if err != nil {
		http.Error(w, "Invalid trainer ID, must be an integer", http.StatusBadRequest)
		return
	}

	availableAppointments, err := services.GetAvailableAppointments(r.Context(), trainerIDInt, startsAt, endsAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting available appointments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(availableAppointments)
}

/*
CreateAppointment POST request to create a new appointment.

Path: /appointments/create

Parameters:
  - trainer_id (int): ID of the trainer for the appointment
  - user_id (int): ID of the user booking the appointment
  - starts_at (string): Start time of appointment in format "2019-01-25T09:00:00-08:00"
  - ends_at (string): End time of appointment in format "2019-01-25T09:00:00-08:00"

Returns:
  - 201 Created on success with the created appointment
  - 400 Bad Request if request body is invalid
*/
func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	appointment := models.Appointment{}

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the required fields
	if appointment.TrainerID == 0 {
		http.Error(w, "Missing trainer_id", http.StatusBadRequest)
		return
	}
	if appointment.UserID == 0 {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}
	if appointment.StartsAt == "" {
		http.Error(w, "Missing starts_at", http.StatusBadRequest)
		return
	}
	if appointment.EndsAt == "" {
		http.Error(w, "Missing ends_at", http.StatusBadRequest)
		return
	}

	// Create the appointment
	err := services.CreateAppointment(r.Context(), appointment.TrainerID, appointment.UserID, appointment.StartsAt, appointment.EndsAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating appointment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

/*
ScheduledAppointments GET request to get a list of scheduled appointments for a trainer.

Path: /appointments/scheduled

Parameters:
  - trainer_id (int): ID of the trainer for the appointment

Returns:
  - 200 OK with list of scheduled appointments
    {
    "id": 1,
    "trainer_id": 1,
    "user_id": 1,
    "starts_at": "2019-01-25T09:00:00-08:00",
    "ends_at": "2019-01-25T09:00:00-08:00"
    }
  - 400 Bad Request if request body is invalid
*/
func ScheduledAppointments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	trainerID := query.Get("trainer_id")

	if trainerID == "" {
		http.Error(w, "Missing trainer_id", http.StatusBadRequest)
		return
	}

	trainerIDInt, err := strconv.Atoi(trainerID)
	if err != nil {
		http.Error(w, "Invalid trainer ID, must be an integer", http.StatusBadRequest)
		return
	}

	appointments, err := services.ScheduledAppointments(r.Context(), trainerIDInt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting scheduled appointments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}
