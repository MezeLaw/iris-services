package handler

import (
	"context"
	"github.com/MezeLaw/iris-services/internal/models"
	"go.uber.org/zap"
)

type PatientsHandler interface {
	Create(context.Context, models.Patient) (*models.Patient, error)
	Get(ctx context.Context) (*models.Patient, error)
	Update(ctx context.Context) (*models.Patient, error)
	Delete(ctx context.Context) error
}

type PatientsService interface {
	CreatePatient(context.Context, models.Patient) (*models.Patient, error)
	GetPatient(context.Context, string) (*models.Patient, error)
	UpdatePatient(context.Context, models.Patient) error
	DeletePatient(context.Context, string) error
}

type Patients struct {
	Service PatientsService
	Logger  *zap.SugaredLogger
}

func New(service PatientsService, logger *zap.SugaredLogger) PatientsHandler {
	return &Patients{Service: service, Logger: logger}
}

func (p *Patients) Create(ctx context.Context, patient models.Patient) (*models.Patient, error) {
	p.Logger.Infof("Creating patient: %s", patient)
	result, err := p.Service.CreatePatient(ctx, patient)
	if err != nil {
		p.Logger.Errorf("Error creating patient: %s", err)
		return nil, err
	}
	return result, nil
}

func (p *Patients) Get(ctx context.Context) (*models.Patient, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) Update(ctx context.Context) (*models.Patient, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) Delete(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
