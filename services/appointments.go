package services

import (
	"appointment-scheduling/data"
	"appointment-scheduling/models"
	"context"
	"fmt"
	"time"
)

const (
	appointmentDuration = 30 * time.Minute
	businessHoursStart  = 8 * time.Hour
	businessHoursEnd    = 17 * time.Hour
)

// GetAvailableAppointments returns a list of available appointment times for a trainer between two dates
func GetAvailableAppointments(ctx context.Context, trainerID int, startDate string, endDate string) ([]string, error) {
	// Parse the start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	bookedSlots, err := getBookedSlots(trainerID)
	if err != nil {
		return nil, err
	}

	pacificLocation, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, err
	}

	// Calculate available start times
	var availableAppointments []string
	for current := start; current.Before(end) || current.Equal(end); current = current.Add(24 * time.Hour) {
		startTime := time.Date(current.Year(), current.Month(), current.Day(), 8, 0, 0, 0, pacificLocation)
		endTime := time.Date(current.Year(), current.Month(), current.Day(), 17, 0, 0, 0, pacificLocation)
		for current := startTime; current.Before(endTime); current = current.Add(appointmentDuration) {
			if !bookedSlots[current.Format(time.RFC3339)] {
				availableAppointments = append(availableAppointments, current.Format(time.RFC3339))
			}
		}
	}

	return availableAppointments, nil
}

// CreateAppointment creates a new appointment
func CreateAppointment(ctx context.Context, trainerID int, userID int, startTime string, endTime string) error {
	// Parse the start and end times
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return err
	}

	err = isValidAppointment(start, end)
	if err != nil {
		return err
	}

	bookedSlots, err := getBookedSlots(trainerID)
	if err != nil {
		return err
	}

	if bookedSlots[startTime] {
		return fmt.Errorf("appointment slot is already booked")
	}

	appointment := &models.Appointment{
		TrainerID: trainerID,
		UserID:    userID,
		StartsAt:  startTime,
		EndsAt:    endTime,
	}

	err = data.CreateAppointment(appointment)
	if err != nil {
		return err
	}

	return nil
}

// ScheduledAppointments returns a list of scheduled appointments for a trainer
func ScheduledAppointments(ctx context.Context, trainerID int) ([]*models.Appointment, error) {
	appointments, err := data.GetAppointments(trainerID)
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

// getBookedSlots returns a map of booked slots for a trainer
func getBookedSlots(trainerID int) (map[string]bool, error) {
	// Get scheduled appointments
	scheduledAppointments, err := data.GetAppointments(trainerID)
	if err != nil {
		return nil, err
	}

	// Get booked start times
	bookedSlots := make(map[string]bool)
	for _, appointment := range scheduledAppointments {
		startTime, _ := time.Parse(time.RFC3339, appointment.StartsAt)
		bookedSlots[startTime.Format(time.RFC3339)] = true
	}

	return bookedSlots, nil
}

func isValidAppointment(start time.Time, end time.Time) error {
	// Verify the appointment duration is 30 minutes
	if end.Sub(start) != appointmentDuration {
		return fmt.Errorf("appointment must be exactly 30 minutes long")
	}

	// Ensure the appointment is within business hours
	startHour := start.Hour()
	endHour := end.Hour()
	if startHour < int(businessHoursStart.Hours()) || endHour > int(businessHoursEnd.Hours()) {
		return fmt.Errorf("appointment must be within business hours (8 AM to 5 PM)")
	}

	return nil
}
