package models

type Appointment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	TrainerID int    `json:"trainer_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
}
