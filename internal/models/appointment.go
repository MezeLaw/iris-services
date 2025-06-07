package models

type AppointmentStatus string

const (
	AppointmentStatusScheduled  AppointmentStatus = "SCHEDULED"
	AppointmentStatusInProgress AppointmentStatus = "IN_PROGRESS"
	AppointmentStatusCompleted  AppointmentStatus = "COMPLETED"
	AppointmentStatusCancelled  AppointmentStatus = "CANCELLED"
)

type AppointmentRequest struct {
	ID        string                 `json:"id,omitempty"`
	ClientID  string                 `json:"client_id"`
	PatientID string                 `json:"patient_id"`
	DoctorID  string                 `json:"doctor_id"`
	Date      string                 `json:"date"`     // Format: RFC3339
	Duration  int                    `json:"duration"` // Duration in minutes
	Status    AppointmentStatus      `json:"status"`
	Notes     string                 `json:"notes,omitempty"`
	CreatedAt string                 `json:"created_at,omitempty"`
	UpdatedAt string                 `json:"updated_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type Appointment struct {
	ID        string                 `dynamodbav:"id"`
	ClientID  string                 `dynamodbav:"client_id"`
	PatientID string                 `dynamodbav:"patient_id"`
	DoctorID  string                 `dynamodbav:"doctor_id"`
	Date      string                 `dynamodbav:"date"`
	Duration  int                    `dynamodbav:"duration"`
	Status    AppointmentStatus      `dynamodbav:"status"`
	Notes     string                 `dynamodbav:"notes,omitempty"`
	CreatedAt string                 `dynamodbav:"created_at"`
	UpdatedAt string                 `dynamodbav:"updated_at"`
	Metadata  map[string]interface{} `dynamodbav:"metadata,omitempty"`
}

type GetAppointmentRequest struct {
	ID        string `json:"id,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	PatientID string `json:"patient_id,omitempty"`
	DoctorID  string `json:"doctor_id,omitempty"`
}
