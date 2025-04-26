package handler

import (
	"context"
	"github.com/MezeLaw/iris-services/internal/models"
	"go.uber.org/zap"
)

type PatientsHandler interface {
	Create(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	Get(ctx context.Context, getPatient *models.GetPatientRequest) (*models.PatientRequest, error)
	GetAll(ctx context.Context, clientID string) ([]*models.PatientRequest, error)
	Update(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error)
	Delete(ctx context.Context, userID string) error
}

type PatientsService interface {
	CreatePatient(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	GetPatient(context.Context, *models.GetPatientRequest) (*models.PatientRequest, error)
	GetAllPatients(context.Context, string) ([]*models.PatientRequest, error)
	UpdatePatient(context.Context, *models.PatientRequest) error
	DeletePatient(context.Context, string) error
}

type Patients struct {
	Service PatientsService
	Logger  *zap.SugaredLogger
}

func New(service PatientsService, logger *zap.SugaredLogger) PatientsHandler {
	return &Patients{Service: service, Logger: logger}
}

func (p *Patients) Create(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error) {
	p.Logger.Infof("Creating patient: %s", patient)
	result, err := p.Service.CreatePatient(ctx, patient)
	if err != nil {
		p.Logger.Errorf("Error creating patient: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *Patients) Get(ctx context.Context, getRequest *models.GetPatientRequest) (*models.PatientRequest, error) {
	p.Logger.Infof("Getting patient with params: %s", getRequest)
	result, err := p.Service.GetPatient(ctx, getRequest)
	if err != nil {
		p.Logger.Errorf("Error getting patient: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *Patients) GetAll(ctx context.Context, clientID string) ([]*models.PatientRequest, error) {
	p.Logger.Infof("Getting all patients with clientID: %s", clientID)
	result, err := p.Service.GetAllPatients(ctx, clientID)
	if err != nil {
		p.Logger.Errorf("Error getting patients by clientId: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *Patients) Update(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error) {
	p.Logger.Infof("Updating patient: %s", patient)
	err := p.Service.UpdatePatient(ctx, patient)
	if err != nil {
		p.Logger.Errorf("Error updating patient: %s", err)
		return nil, err
	}
	// Return the updated patient data
	return patient, nil
}

func (p *Patients) Delete(ctx context.Context, userID string) error {
	p.Logger.Infof("Deleting patient with userID: %s", userID)
	err := p.Service.DeletePatient(ctx, userID)
	if err != nil {
		p.Logger.Errorf("Error deleting patient: %s", err)
		return err
	}
	return nil
}
