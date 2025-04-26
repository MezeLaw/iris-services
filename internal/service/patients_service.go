package service

import (
	"context"
	"fmt"
	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type PatientsRepository interface {
	Save(ctx context.Context, p *models.Patient) error
	GetByID(ctx context.Context, id string) (*models.Patient, error)
	GetByClientID(ctx context.Context, clientID string) ([]*models.Patient, error)
	GetByDocument(ctx context.Context, docType, docNumber string) (*models.Patient, error)
	Delete(ctx context.Context, id string) error
}

type PatientsService interface {
	CreatePatient(context.Context, *models.PatientRequest) (*models.PatientRequest, error)
	GetPatient(context.Context, *models.GetPatientRequest) (*models.PatientRequest, error)
	GetAllPatients(context.Context, string) ([]*models.PatientRequest, error)
	UpdatePatient(context.Context, *models.PatientRequest) error
	DeletePatient(context.Context, string) error
}

type Patients struct {
	Logger             *zap.SugaredLogger
	PatientsRepository PatientsRepository
}

func New(logger *zap.SugaredLogger, repository PatientsRepository) PatientsService {
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

func (p *Patients) GetPatient(ctx context.Context, params *models.GetPatientRequest) (*models.PatientRequest, error) {
	// Si se proporciona un ID, buscar por ID
	if params.ID != "" {
		p.Logger.Info("Getting patient by ID", zap.String("id", params.ID))
		patient, err := p.PatientsRepository.GetByID(ctx, params.ID)
		if err != nil {
			p.Logger.Error("Error getting patient by ID", zap.String("id", params.ID), zap.Error(err))
			return nil, err
		}
		return p.mapPatientToRequest(patient), nil
	}

	// Si se proporcionan DocType y DocNumber, buscar por documento
	if params.DocType != "" && params.DocNumber != "" {
		p.Logger.Info("Getting patient by Document",
			zap.String("docType", params.DocType),
			zap.String("docNumber", params.DocNumber))

		patient, err := p.PatientsRepository.GetByDocument(ctx, params.DocType, params.DocNumber)
		if err != nil {
			p.Logger.Error("Error getting patient by Document",
				zap.String("docType", params.DocType),
				zap.String("docNumber", params.DocNumber),
				zap.Error(err))
			return nil, err
		}

		return p.mapPatientToRequest(patient), nil
	}

	// Si no se proporcionó ningún parámetro válido para la búsqueda
	p.Logger.Error("Invalid parameters for GetPatient")
	return nil, fmt.Errorf("invalid parameters: must provide ID, ClientID, or DocType/DocNumber")
}

func (p *Patients) GetAllPatients(ctx context.Context, identifier string) ([]*models.PatientRequest, error) {
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
		UpdatedAt:      time.Now().Format(time.RFC3339),
		Metadata:       req.Metadata,
	}
}

func (p *Patients) mapPatientToRequest(patient *models.Patient) *models.PatientRequest {
	return &models.PatientRequest{
		ClientID:       patient.ClientID,
		FirstName:      patient.FirstName,
		LastName:       patient.LastName,
		DocType:        patient.DocType,
		DocNumber:      patient.DocNumber,
		BirthDate:      patient.BirthDate,
		Gender:         patient.Gender,
		CountryCode:    patient.CountryCode,
		PhoneNumber:    patient.PhoneNumber,
		Email:          patient.Email,
		AddressStreet:  patient.AddressStreet,
		AddressNumber:  patient.AddressNumber,
		AddressCity:    patient.AddressCity,
		AddressCountry: patient.AddressCountry,
		ZipCode:        patient.ZipCode,
		CreatedAt:      patient.CreatedAt,
		UpdatedAt:      patient.UpdatedAt,
		Metadata:       patient.Metadata,
	}
}
