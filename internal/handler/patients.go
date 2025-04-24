package handler

import (
	"context"
	"github.com/MezeLaw/iris-services/internal/models"
	"go.uber.org/zap"
)

type PatientsHandler interface {
	Create(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	Get(ctx context.Context, userID string) (*models.PatientRequest, error)
	Update(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error)
	Delete(ctx context.Context, userID string) error
}

type PatientsService interface {
	CreatePatient(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	GetPatient(context.Context, string) (*models.PatientRequest, error)
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

func (p *Patients) Get(ctx context.Context, userID string) (*models.PatientRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) Update(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) Delete(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}
