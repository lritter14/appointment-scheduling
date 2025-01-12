package data

import (
	"appointment-scheduling/models"
	"errors"
)

var insertAppointment = `
	INSERT INTO APPOINTMENTS (id, trainer_id, user_id, starts_at, ends_at)
	VALUES (?, ?, ?, ?, ?)
`

var getAppointments = `
	SELECT * FROM APPOINTMENTS where trainer_id = ?
`

// CreateAppointment creates a new appointment in the database
func CreateAppointment(appointment *models.Appointment) error {
	if appointment == nil {
		return errors.New("appointment is nil")
	}
	if appointment.ID == 0 {
		// Get the highest appointment ID and increment by 1
		var maxID int
		err := db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM APPOINTMENTS").Scan(&maxID)
		if err != nil {
			return err
		}
		appointment.ID = maxID + 1
	}
	_, err := db.Exec(insertAppointment, appointment.ID, appointment.TrainerID, appointment.UserID, appointment.StartsAt, appointment.EndsAt)
	if err != nil {
		return err
	}
	return nil
}

// GetAppointments returns all appointments for a trainer
func GetAppointments(trainerID int) ([]*models.Appointment, error) {
	rows, err := db.Query(getAppointments, trainerID)
	if err != nil {
		return nil, err
	}

	var appointments []*models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		err := rows.Scan(&appointment.ID, &appointment.TrainerID, &appointment.UserID, &appointment.StartsAt, &appointment.EndsAt)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}
