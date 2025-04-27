package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockPatientsRepository is a mock implementation of PatientsRepository
type MockPatientsRepository struct {
	mock.Mock
}

func (m *MockPatientsRepository) Save(ctx context.Context, p *models.Patient) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPatientsRepository) GetByID(ctx context.Context, id string) (*models.Patient, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientsRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Patient, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Patient), args.Error(1)
}

func (m *MockPatientsRepository) GetByDocument(ctx context.Context, docType, docNumber string) (*models.Patient, error) {
	args := m.Called(ctx, docType, docNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Patient), args.Error(1)
}

func (m *MockPatientsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test setup helper function
func setupTest() (*Patients, *MockPatientsRepository) {
	mockRepo := new(MockPatientsRepository)
	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()
	service := &Patients{
		Logger:             sugarLogger,
		PatientsRepository: mockRepo,
	}
	
	return service, mockRepo
}

// Sample patient data helper function
func createSamplePatientRequest() *models.PatientRequest {
	return &models.PatientRequest{
		ClientID:       "client123",
		FirstName:      "John",
		LastName:       "Doe",
		DocType:        "DNI",
		DocNumber:      "12345678",
		BirthDate:      "1990-01-01",
		Gender:         "M",
		CountryCode:    "54",
		PhoneNumber:    "1234567890",
		Email:          "john.doe@example.com",
		AddressStreet:  "Main St",
		AddressNumber:  "123",
		AddressCity:    "Buenos Aires",
		AddressCountry: "Argentina",
		ZipCode:        "1234",
		Metadata:       map[string]interface{}{"key": "value"},
	}
}

func createSamplePatient(id string) *models.Patient {
	return &models.Patient{
		ID:             id,
		ClientID:       "client123",
		FirstName:      "John",
		LastName:       "Doe",
		DocType:        "DNI",
		DocNumber:      "12345678",
		DocKey:         "DNI#12345678",
		BirthDate:      "1990-01-01",
		Gender:         "M",
		CountryCode:    "54",
		PhoneNumber:    "1234567890",
		Email:          "john.doe@example.com",
		AddressStreet:  "Main St",
		AddressNumber:  "123",
		AddressCity:    "Buenos Aires",
		AddressCountry: "Argentina",
		ZipCode:        "1234",
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
		Metadata:       map[string]interface{}{"key": "value"},
	}
}

// Tests for CreatePatient
func TestPatients_CreatePatient_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	req := createSamplePatientRequest()
	
	// Expect the repository to be called with any patient object
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Patient")).Return(nil)
	
	// Execute
	result, err := service.CreatePatient(ctx, req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, req, result)
	mockRepo.AssertExpectations(t)
}

func TestPatients_CreatePatient_Error(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	req := createSamplePatientRequest()
	expectedErr := errors.New("database error")
	
	// Expect the repository to return an error
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Patient")).Return(expectedErr)
	
	// Execute
	result, err := service.CreatePatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests for GetPatient
func TestPatients_GetPatient_ByID_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	req := &models.GetPatientRequest{ID: patientID}
	
	expectedPatient := createSamplePatient(patientID)
	mockRepo.On("GetByID", ctx, patientID).Return(expectedPatient, nil)
	
	// Execute
	result, err := service.GetPatient(ctx, req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, patientID, result.ID)
	assert.Equal(t, expectedPatient.FirstName, result.FirstName)
	mockRepo.AssertExpectations(t)
}

func TestPatients_GetPatient_ByID_Error(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	req := &models.GetPatientRequest{ID: patientID}
	expectedErr := errors.New("not found")
	
	mockRepo.On("GetByID", ctx, patientID).Return(nil, expectedErr)
	
	// Execute
	result, err := service.GetPatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPatients_GetPatient_ByDocument_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	docType := "DNI"
	docNumber := "12345678"
	req := &models.GetPatientRequest{DocType: docType, DocNumber: docNumber}
	
	expectedPatient := createSamplePatient("patient123")
	mockRepo.On("GetByDocument", ctx, docType, docNumber).Return(expectedPatient, nil)
	
	// Execute
	result, err := service.GetPatient(ctx, req)
	
	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPatient.ID, result.ID)
	assert.Equal(t, docType, result.DocType)
	assert.Equal(t, docNumber, result.DocNumber)
	mockRepo.AssertExpectations(t)
}

func TestPatients_GetPatient_ByDocument_Error(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	docType := "DNI"
	docNumber := "12345678"
	req := &models.GetPatientRequest{DocType: docType, DocNumber: docNumber}
	expectedErr := errors.New("not found")
	
	mockRepo.On("GetByDocument", ctx, docType, docNumber).Return(nil, expectedErr)
	
	// Execute
	result, err := service.GetPatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPatients_GetPatient_InvalidParams(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()
	req := &models.GetPatientRequest{} // Empty request
	
	// Execute
	result, err := service.GetPatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid parameters")
}

// Tests for GetAllPatients
func TestPatients_GetAllPatients_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	clientID := "client123"
	
	patient1 := createSamplePatient("patient1")
	patient2 := createSamplePatient("patient2")
	patients := []*models.Patient{patient1, patient2}
	
	mockRepo.On("GetByClientID", ctx, clientID).Return(patients, nil)
	
	// Execute
	result, err := service.GetAllPatients(ctx, clientID)
	
	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, patient1.ID, result[0].ID)
	assert.Equal(t, patient2.ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestPatients_GetAllPatients_EmptyClientID(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()
	
	// Execute
	result, err := service.GetAllPatients(ctx, "")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "client-id cannot be empty")
}

func TestPatients_GetAllPatients_RepositoryError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	clientID := "client123"
	expectedErr := errors.New("database error")
	
	mockRepo.On("GetByClientID", ctx, clientID).Return(nil, expectedErr)
	
	// Execute
	result, err := service.GetAllPatients(ctx, clientID)
	
	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests for UpdatePatient
func TestPatients_UpdatePatient_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	
	req := createSamplePatientRequest()
	req.ID = patientID
	
	existingPatient := createSamplePatient(patientID)
	
	mockRepo.On("GetByID", ctx, patientID).Return(existingPatient, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Patient")).Return(nil)
	
	// Execute
	err := service.UpdatePatient(ctx, req)
	
	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatients_UpdatePatient_NoID(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()
	req := createSamplePatientRequest()
	// ID intentionally left empty
	
	// Execute
	err := service.UpdatePatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "patient ID is required")
}

func TestPatients_UpdatePatient_NotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	req := createSamplePatientRequest()
	req.ID = patientID
	
	expectedErr := errors.New("not found")
	mockRepo.On("GetByID", ctx, patientID).Return(nil, expectedErr)
	
	// Execute
	err := service.UpdatePatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find patient")
	mockRepo.AssertExpectations(t)
}

func TestPatients_UpdatePatient_SaveError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	req := createSamplePatientRequest()
	req.ID = patientID
	
	existingPatient := createSamplePatient(patientID)
	expectedErr := errors.New("database error")
	
	mockRepo.On("GetByID", ctx, patientID).Return(existingPatient, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Patient")).Return(expectedErr)
	
	// Execute
	err := service.UpdatePatient(ctx, req)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update patient")
	mockRepo.AssertExpectations(t)
}

// Tests for DeletePatient
func TestPatients_DeletePatient_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	
	existingPatient := createSamplePatient(patientID)
	
	mockRepo.On("GetByID", ctx, patientID).Return(existingPatient, nil)
	mockRepo.On("Delete", ctx, patientID).Return(nil)
	
	// Execute
	err := service.DeletePatient(ctx, patientID)
	
	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPatients_DeletePatient_EmptyID(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()
	
	// Execute
	err := service.DeletePatient(ctx, "")
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "patient ID cannot be empty")
}

func TestPatients_DeletePatient_NotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	
	expectedErr := errors.New("not found")
	mockRepo.On("GetByID", ctx, patientID).Return(nil, expectedErr)
	
	// Execute
	err := service.DeletePatient(ctx, patientID)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find patient")
	mockRepo.AssertExpectations(t)
}

func TestPatients_DeletePatient_DeleteError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	
	existingPatient := createSamplePatient(patientID)
	expectedErr := errors.New("database error")
	
	mockRepo.On("GetByID", ctx, patientID).Return(existingPatient, nil)
	mockRepo.On("Delete", ctx, patientID).Return(expectedErr)
	
	// Execute
	err := service.DeletePatient(ctx, patientID)
	
	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete patient")
	mockRepo.AssertExpectations(t)
}

// Tests for mapper methods
func TestPatients_MapRequestToPatient(t *testing.T) {
	// Setup
	service, _ := setupTest()
	req := createSamplePatientRequest()
	
	// Execute
	patient := service.mapRequestToPatient(req)
	
	// Assert
	assert.NotEmpty(t, patient.ID) // ID should be generated
	assert.Equal(t, req.ClientID, patient.ClientID)
	assert.Equal(t, req.FirstName, patient.FirstName)
	assert.Equal(t, req.LastName, patient.LastName)
	assert.Equal(t, req.DocType, patient.DocType)
	assert.Equal(t, req.DocNumber, patient.DocNumber)
	assert.Equal(t, req.Metadata, patient.Metadata)
	assert.NotEmpty(t, patient.CreatedAt)
	assert.NotEmpty(t, patient.UpdatedAt)
}

func TestPatients_MapPatientToRequest(t *testing.T) {
	// Setup
	service, _ := setupTest()
	patientID := "test-id-123"
	patient := createSamplePatient(patientID)
	
	// Execute
	req := service.mapPatientToRequest(patient)
	
	// Assert
	assert.Equal(t, patient.ID, req.ID)
	assert.Equal(t, patient.ClientID, req.ClientID)
	assert.Equal(t, patient.FirstName, req.FirstName)
	assert.Equal(t, patient.LastName, req.LastName)
	assert.Equal(t, patient.DocType, req.DocType)
	assert.Equal(t, patient.DocNumber, req.DocNumber)
	assert.Equal(t, patient.CreatedAt, req.CreatedAt)
	assert.Equal(t, patient.UpdatedAt, req.UpdatedAt)
	assert.Equal(t, patient.Metadata, req.Metadata)
}

// Test for New
func TestNew(t *testing.T) {
	// Setup
	mockRepo := new(MockPatientsRepository)
	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()
	
	// Execute
	service := New(sugarLogger, mockRepo)
	
	// Assert
	assert.NotNil(t, service)
	assert.Implements(t, (*PatientsService)(nil), service)
}
