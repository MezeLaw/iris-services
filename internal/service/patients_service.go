package service

import (
	"context"
	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type PatientsRepository interface {
	Save(ctx context.Context, p *models.Patient) error
	Get(ctx context.Context, id string) (*models.Patient, error)
	Update(ctx context.Context, p *models.Patient) error
	Delete(ctx context.Context, id string) error
}

type PatientsService interface {
	CreatePatient(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	GetPatient(context.Context, string) (*models.PatientRequest, error)
	UpdatePatient(context.Context, *models.PatientRequest) error
	DeletePatient(context.Context, string) error
}

type Patients struct {
	Logger             *zap.Logger
	PatientsRepository PatientsRepository
}

func New(logger *zap.Logger, repository PatientsRepository) PatientsService {
	return &Patients{
		Logger:             logger,
		PatientsRepository: repository,
	}
}

func (p *Patients) CreatePatient(ctx context.Context, request *models.PatientRequest) (*models.PatientRequest, error) {
	patient := p.mapRequestToPatient(request)
	if err := p.PatientsRepository.Save(ctx, patient); err != nil {
		p.Logger.Error("Error on PatientsRepository.Save", zap.Error(err))
		return nil, err
	}
	return request, nil
}

func (p *Patients) GetPatient(ctx context.Context, s string) (*models.PatientRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) UpdatePatient(ctx context.Context, request *models.PatientRequest) error {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) DeletePatient(ctx context.Context, s string) error {
	//TODO implement me
	panic("implement me")
}

func (p *Patients) mapRequestToPatient(req *models.PatientRequest) *models.Patient {
	return &models.Patient{
		ID:             uuid.NewString(),
		ClientID:       req.ClientID,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DocType:        req.DocType,
		DocNumber:      req.DocNumber,
		BirthDate:      req.BirthDate,
		Gender:         req.Gender,
		CountryCode:    req.CountryCode,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		AddressStreet:  req.AddressStreet,
		AddressNumber:  req.AddressNumber,
		AddressCity:    req.AddressCity,
		AddressCountry: req.AddressCountry,
		ZipCode:        req.ZipCode,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
}
