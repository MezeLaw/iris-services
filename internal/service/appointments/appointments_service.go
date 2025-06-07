package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AppointmentsRepository interface {
	Save(ctx context.Context, a *models.Appointment) error
	GetByID(ctx context.Context, id string) (*models.Appointment, error)
	GetByClientID(ctx context.Context, clientID string) ([]*models.Appointment, error)
	GetByPatientID(ctx context.Context, patientID string) ([]*models.Appointment, error)
	GetByDoctorID(ctx context.Context, doctorID string) ([]*models.Appointment, error)
	Delete(ctx context.Context, id string) error
}

type AppointmentsService interface {
	CreateAppointment(context.Context, *models.AppointmentRequest) (*models.AppointmentRequest, error)
	GetAppointment(context.Context, *models.GetAppointmentRequest) (*models.AppointmentRequest, error)
	GetAllAppointments(context.Context, string) ([]*models.AppointmentRequest, error)
	UpdateAppointment(context.Context, *models.AppointmentRequest) error
	DeleteAppointment(context.Context, string) error
}

type Appointments struct {
	Logger                 *zap.SugaredLogger
	AppointmentsRepository AppointmentsRepository
}

func New(logger *zap.SugaredLogger, repository AppointmentsRepository) AppointmentsService {
	return &Appointments{
		Logger:                 logger,
		AppointmentsRepository: repository,
	}
}

func (a *Appointments) CreateAppointment(ctx context.Context, request *models.AppointmentRequest) (*models.AppointmentRequest, error) {
	// Validar status
	if err := validateStatus(request.Status); err != nil {
		a.Logger.Error("Invalid status value", zap.String("status", string(request.Status)))
		return nil, err
	}

	appointment := a.mapRequestToAppointment(request)
	if err := a.AppointmentsRepository.Save(ctx, appointment); err != nil {
		a.Logger.Error("Error on AppointmentsRepository.Save", zap.Error(err))
		return nil, err
	}
	return request, nil
}

func (a *Appointments) GetAppointment(ctx context.Context, params *models.GetAppointmentRequest) (*models.AppointmentRequest, error) {
	// Si se proporciona un ID, buscar por ID
	if params.ID != "" {
		a.Logger.Info("Getting appointment by ID", zap.String("id", params.ID))
		appointment, err := a.AppointmentsRepository.GetByID(ctx, params.ID)
		if err != nil {
			a.Logger.Error("Error getting appointment by ID", zap.String("id", params.ID), zap.Error(err))
			return nil, err
		}
		return a.mapAppointmentToRequest(appointment), nil
	}

	// Si se proporciona un PatientID, buscar por PatientID
	if params.PatientID != "" {
		a.Logger.Info("Getting appointments by PatientID", zap.String("patientID", params.PatientID))
		appointments, err := a.AppointmentsRepository.GetByPatientID(ctx, params.PatientID)
		if err != nil {
			a.Logger.Error("Error getting appointments by PatientID", zap.String("patientID", params.PatientID), zap.Error(err))
			return nil, err
		}
		if len(appointments) == 0 {
			return nil, fmt.Errorf("no appointments found for patientID: %s", params.PatientID)
		}
		return a.mapAppointmentToRequest(appointments[0]), nil
	}

	// Si se proporciona un DoctorID, buscar por DoctorID
	if params.DoctorID != "" {
		a.Logger.Info("Getting appointments by DoctorID", zap.String("doctorID", params.DoctorID))
		appointments, err := a.AppointmentsRepository.GetByDoctorID(ctx, params.DoctorID)
		if err != nil {
			a.Logger.Error("Error getting appointments by DoctorID", zap.String("doctorID", params.DoctorID), zap.Error(err))
			return nil, err
		}
		if len(appointments) == 0 {
			return nil, fmt.Errorf("no appointments found for doctorID: %s", params.DoctorID)
		}
		return a.mapAppointmentToRequest(appointments[0]), nil
	}

	// Si no se proporcionó ningún parámetro válido para la búsqueda
	a.Logger.Error("Invalid parameters for GetAppointment")
	return nil, fmt.Errorf("invalid parameters: must provide ID, ClientID, PatientID, or DoctorID")
}

func (a *Appointments) GetAllAppointments(ctx context.Context, identifier string) ([]*models.AppointmentRequest, error) {
	// Verificar que el identificador no esté vacío
	if identifier == "" {
		a.Logger.Error("Error: empty client-id provided to GetAllAppointments")
		return nil, fmt.Errorf("client-id cannot be empty")
	}

	a.Logger.Info("Getting all appointments by ClientID", zap.String("clientID", identifier))

	// Obtener citas del repositorio usando el cliente ID
	appointments, err := a.AppointmentsRepository.GetByClientID(ctx, identifier)
	if err != nil {
		a.Logger.Error("Error getting appointments by ClientID", zap.String("clientID", identifier), zap.Error(err))
		return nil, err
	}

	// Mapear las citas del modelo de BD al modelo de request
	appointmentRequests := make([]*models.AppointmentRequest, 0, len(appointments))
	for _, appointment := range appointments {
		appointmentRequests = append(appointmentRequests, a.mapAppointmentToRequest(appointment))
	}

	a.Logger.Info("Successfully retrieved appointments", zap.Int("count", len(appointmentRequests)))
	return appointmentRequests, nil
}

func (a *Appointments) UpdateAppointment(ctx context.Context, request *models.AppointmentRequest) error {
	// Verificar que la cita tenga un ID
	if request.ID == "" {
		a.Logger.Error("Error: Missing appointment ID for update")
		return fmt.Errorf("appointment ID is required for update")
	}

	// Validar status
	if err := validateStatus(request.Status); err != nil {
		a.Logger.Error("Invalid status value", zap.String("status", string(request.Status)))
		return err
	}

	a.Logger.Info("Updating appointment", zap.String("id", request.ID))

	// Obtener la cita existente
	existingAppointment, err := a.AppointmentsRepository.GetByID(ctx, request.ID)
	if err != nil {
		a.Logger.Error("Error fetching appointment to update", zap.String("id", request.ID), zap.Error(err))
		return fmt.Errorf("failed to find appointment with ID %s: %w", request.ID, err)
	}

	// Actualizar los campos de la cita existente
	updatedAppointment := &models.Appointment{
		ID:        existingAppointment.ID,
		ClientID:  request.ClientID,
		PatientID: request.PatientID,
		DoctorID:  request.DoctorID,
		Date:      request.Date,
		Duration:  request.Duration,
		Status:    request.Status,
		Notes:     request.Notes,
		CreatedAt: existingAppointment.CreatedAt,
		UpdatedAt: time.Now().Format(time.RFC3339),
		Metadata:  request.Metadata,
	}

	// Guardar la cita actualizada
	if err := a.AppointmentsRepository.Save(ctx, updatedAppointment); err != nil {
		a.Logger.Error("Error updating appointment", zap.String("id", request.ID), zap.Error(err))
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	a.Logger.Info("Appointment updated successfully", zap.String("id", request.ID))
	return nil
}

func (a *Appointments) DeleteAppointment(ctx context.Context, id string) error {
	// Verificar que el ID no esté vacío
	if id == "" {
		a.Logger.Error("Error: empty ID provided for appointment deletion")
		return fmt.Errorf("appointment ID cannot be empty")
	}

	a.Logger.Info("Deleting appointment", zap.String("id", id))

	// Verificar primero si la cita existe
	_, err := a.AppointmentsRepository.GetByID(ctx, id)
	if err != nil {
		a.Logger.Error("Error finding appointment to delete", zap.String("id", id), zap.Error(err))
		return fmt.Errorf("failed to find appointment with ID %s: %w", id, err)
	}

	// Eliminar la cita
	if err := a.AppointmentsRepository.Delete(ctx, id); err != nil {
		a.Logger.Error("Error deleting appointment", zap.String("id", id), zap.Error(err))
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	a.Logger.Info("Appointment deleted successfully", zap.String("id", id))
	return nil
}

func (a *Appointments) mapRequestToAppointment(req *models.AppointmentRequest) *models.Appointment {
	return &models.Appointment{
		ID:        uuid.NewString(),
		ClientID:  req.ClientID,
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		Date:      req.Date,
		Duration:  req.Duration,
		Status:    req.Status,
		Notes:     req.Notes,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Metadata:  req.Metadata,
	}
}

func (a *Appointments) mapAppointmentToRequest(appointment *models.Appointment) *models.AppointmentRequest {
	return &models.AppointmentRequest{
		ID:        appointment.ID,
		ClientID:  appointment.ClientID,
		PatientID: appointment.PatientID,
		DoctorID:  appointment.DoctorID,
		Date:      appointment.Date,
		Duration:  appointment.Duration,
		Status:    appointment.Status,
		Notes:     appointment.Notes,
		CreatedAt: appointment.CreatedAt,
		UpdatedAt: appointment.UpdatedAt,
		Metadata:  appointment.Metadata,
	}
}

func validateStatus(status models.AppointmentStatus) error {
	switch status {
	case models.AppointmentStatusScheduled, models.AppointmentStatusInProgress, models.AppointmentStatusCompleted, models.AppointmentStatusCancelled:
		return nil
	default:
		return fmt.Errorf("invalid status value: %s. Must be one of: %s, %s, %s, %s",
			status,
			models.AppointmentStatusScheduled,
			models.AppointmentStatusInProgress,
			models.AppointmentStatusCompleted,
			models.AppointmentStatusCancelled)
	}
}
