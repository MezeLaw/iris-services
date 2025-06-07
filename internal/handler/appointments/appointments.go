package handler

import (
	"context"
	"fmt"

	"github.com/MezeLaw/iris-services/internal/models"
	"go.uber.org/zap"
)

type AppointmentsHandler interface {
	Create(context.Context, *models.AppointmentRequest) (*models.AppointmentRequest, error)
	Get(ctx context.Context, getAppointment *models.GetAppointmentRequest) (*models.AppointmentRequest, error)
	GetAll(ctx context.Context, clientID string) ([]*models.AppointmentRequest, error)
	Update(ctx context.Context, appointment *models.AppointmentRequest) (*models.AppointmentRequest, error)
	Delete(ctx context.Context, appointmentID string) error
}

type AppointmentsService interface {
	CreateAppointment(context.Context, *models.AppointmentRequest) (*models.AppointmentRequest, error)
	GetAppointment(context.Context, *models.GetAppointmentRequest) (*models.AppointmentRequest, error)
	GetAllAppointments(context.Context, string) ([]*models.AppointmentRequest, error)
	UpdateAppointment(context.Context, *models.AppointmentRequest) error
	DeleteAppointment(context.Context, string) error
}

type Appointments struct {
	Service AppointmentsService
	Logger  *zap.SugaredLogger
}

func New(service AppointmentsService, logger *zap.SugaredLogger) AppointmentsHandler {
	return &Appointments{Service: service, Logger: logger}
}

func (a *Appointments) Create(ctx context.Context, appointment *models.AppointmentRequest) (*models.AppointmentRequest, error) {
	a.Logger.Infof("Creating appointment: %s", appointment)
	if appointment.Status != models.AppointmentStatusScheduled &&
		appointment.Status != models.AppointmentStatusInProgress &&
		appointment.Status != models.AppointmentStatusCompleted &&
		appointment.Status != models.AppointmentStatusCancelled {
		err := fmt.Errorf("invalid status value: %s. Must be one of: %s, %s, %s, %s",
			appointment.Status,
			models.AppointmentStatusScheduled,
			models.AppointmentStatusInProgress,
			models.AppointmentStatusCompleted,
			models.AppointmentStatusCancelled)
		a.Logger.Error(err)
		return nil, err
	}
	result, err := a.Service.CreateAppointment(ctx, appointment)
	if err != nil {
		a.Logger.Errorf("Error creating appointment: %s", err)
		return nil, err
	}
	return result, nil
}

func (a *Appointments) Get(ctx context.Context, getRequest *models.GetAppointmentRequest) (*models.AppointmentRequest, error) {
	a.Logger.Infof("Getting appointment with params: %s", getRequest)
	result, err := a.Service.GetAppointment(ctx, getRequest)
	if err != nil {
		a.Logger.Errorf("Error getting appointment: %s", err)
		return nil, err
	}
	return result, nil
}

func (a *Appointments) GetAll(ctx context.Context, clientID string) ([]*models.AppointmentRequest, error) {
	a.Logger.Infof("Getting all appointments with clientID: %s", clientID)
	result, err := a.Service.GetAllAppointments(ctx, clientID)
	if err != nil {
		a.Logger.Errorf("Error getting appointments by clientId: %s", err)
		return nil, err
	}
	return result, nil
}

func (a *Appointments) Update(ctx context.Context, appointment *models.AppointmentRequest) (*models.AppointmentRequest, error) {
	a.Logger.Infof("Updating appointment: %s", appointment)
	if appointment.Status != models.AppointmentStatusScheduled &&
		appointment.Status != models.AppointmentStatusInProgress &&
		appointment.Status != models.AppointmentStatusCompleted &&
		appointment.Status != models.AppointmentStatusCancelled {
		err := fmt.Errorf("invalid status value: %s. Must be one of: %s, %s, %s, %s",
			appointment.Status,
			models.AppointmentStatusScheduled,
			models.AppointmentStatusInProgress,
			models.AppointmentStatusCompleted,
			models.AppointmentStatusCancelled)
		a.Logger.Error(err)
		return nil, err
	}
	err := a.Service.UpdateAppointment(ctx, appointment)
	if err != nil {
		a.Logger.Errorf("Error updating appointment: %s", err)
		return nil, err
	}
	// Return the updated appointment data
	return appointment, nil
}

func (a *Appointments) Delete(ctx context.Context, appointmentID string) error {
	a.Logger.Infof("Deleting appointment with appointmentID: %s", appointmentID)
	err := a.Service.DeleteAppointment(ctx, appointmentID)
	if err != nil {
		a.Logger.Errorf("Error deleting appointment: %s", err)
		return err
	}
	return nil
}
