package main

import (
	"appointment-scheduling/models"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db               *sql.DB
	appointmentsFile = "data/appointments.json"
)

func createTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS APPOINTMENTS (
			id INTEGER PRIMARY KEY,
			trainer_id INTEGER,
			user_id INTEGER,
			starts_at DATETIME,
			ends_at DATETIME
		)
	`
	_, err := db.Exec(query)
	return err
}

func loadAppointments(appointmentsFile string) error {
	body, err := os.ReadFile(appointmentsFile)
	if err != nil {
		return err
	}

	var appointments []*models.Appointment
	err = json.Unmarshal(body, &appointments)
	if err != nil {
		return err
	}

	for _, appointment := range appointments {

		err = createAppointment(appointment)
		if err != nil {
			return err
		}
	}

	return nil
}

func createAppointment(appointment *models.Appointment) error {

	query := `
		INSERT INTO APPOINTMENTS (id, trainer_id, user_id, starts_at, ends_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, appointment.ID, appointment.TrainerID, appointment.UserID, appointment.StartsAt, appointment.EndsAt)
	return err
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./appointments.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = createTables()
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	err = loadAppointments(appointmentsFile)
	if err != nil {
		log.Fatalf("Failed to load appointments: %v", err)
	}

	log.Println("Appointments loaded successfully.")
}
