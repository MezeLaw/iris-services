package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// MockPatientsService implementa la interfaz PatientsService para los tests
type MockPatientsService struct {
	mock.Mock
}

func (m *MockPatientsService) CreatePatient(ctx context.Context, patient *models.PatientRequest) (*models.PatientRequest, error) {
	args := m.Called(ctx, patient)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientRequest), args.Error(1)
}

func (m *MockPatientsService) GetPatient(ctx context.Context, req *models.GetPatientRequest) (*models.PatientRequest, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PatientRequest), args.Error(1)
}

func (m *MockPatientsService) GetAllPatients(ctx context.Context, clientID string) ([]*models.PatientRequest, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PatientRequest), args.Error(1)
}

func (m *MockPatientsService) UpdatePatient(ctx context.Context, patient *models.PatientRequest) error {
	args := m.Called(ctx, patient)
	return args.Error(0)
}

func (m *MockPatientsService) DeletePatient(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t).Sugar()
	mockService := new(MockPatientsService)

	// Act
	handler := New(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.IsType(t, &Patients{}, handler)
}

func TestPatients_Create(t *testing.T) {
	// Sample patient for tests
	samplePatient := &models.PatientRequest{
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
		AddressCity:    "City",
		AddressCountry: "Country",
		ZipCode:        "12345",
		Metadata:       map[string]interface{}{"key": "value"},
	}

	errorPatient := &models.PatientRequest{
		ClientID:       "client456",
		FirstName:      "Jane",
		LastName:       "Smith",
		DocType:        "DNI",
		DocNumber:      "87654321",
		BirthDate:      "1995-05-05",
		Gender:         "F",
		CountryCode:    "54",
		PhoneNumber:    "9876543210",
		Email:          "jane.smith@example.com",
		AddressStreet:  "Second St",
		AddressNumber:  "456",
		AddressCity:    "Another City",
		AddressCountry: "Another Country",
		ZipCode:        "54321",
		Metadata:       map[string]interface{}{"key2": "value2"},
	}

	// Test scenarios
	tests := []struct {
		name           string
		patient        *models.PatientRequest
		mockSetup      func(*MockPatientsService, *models.PatientRequest)
		expectedError  error
		expectedResult *models.PatientRequest
	}{
		{
			name:    "Success",
			patient: samplePatient,
			mockSetup: func(m *MockPatientsService, p *models.PatientRequest) {
				m.On("CreatePatient", mock.Anything, p).Return(p, nil)
			},
			expectedError:  nil,
			expectedResult: samplePatient,
		},
		{
			name:    "Service Error",
			patient: errorPatient,
			mockSetup: func(m *MockPatientsService, p *models.PatientRequest) {
				m.On("CreatePatient", mock.Anything, p).Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
		{
			name: "Success with non-binary gender",
			patient: &models.PatientRequest{
				ClientID:       "client789",
				FirstName:      "Alex",
				LastName:       "Taylor",
				DocType:        "DNI",
				DocNumber:      "98765432",
				BirthDate:      "1992-03-15",
				Gender:         "NB",
				CountryCode:    "54",
				PhoneNumber:    "5555555555",
				Email:          "alex.taylor@example.com",
				AddressStreet:  "Third St",
				AddressNumber:  "789",
				AddressCity:    "Another City",
				AddressCountry: "Country",
				ZipCode:        "67890",
				Metadata:       map[string]interface{}{"key3": "value3"},
			},
			mockSetup: func(m *MockPatientsService, p *models.PatientRequest) {
				m.On("CreatePatient", mock.Anything, p).Return(p, nil)
			},
			expectedError: nil,
			expectedResult: &models.PatientRequest{
				ClientID:       "client789",
				FirstName:      "Alex",
				LastName:       "Taylor",
				DocType:        "DNI",
				DocNumber:      "98765432",
				BirthDate:      "1992-03-15",
				Gender:         "NB",
				CountryCode:    "54",
				PhoneNumber:    "5555555555",
				Email:          "alex.taylor@example.com",
				AddressStreet:  "Third St",
				AddressNumber:  "789",
				AddressCity:    "Another City",
				AddressCountry: "Country",
				ZipCode:        "67890",
				Metadata:       map[string]interface{}{"key3": "value3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockPatientsService)
			tt.mockSetup(mockService, tt.patient)

			handler := &Patients{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.Create(context.Background(), tt.patient)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestPatients_Get(t *testing.T) {
	// Sample patient for tests
	samplePatient := &models.PatientRequest{
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
		AddressCity:    "City",
		AddressCountry: "Country",
		ZipCode:        "12345",
		Metadata:       map[string]interface{}{"key": "value"},
	}

	// Test scenarios
	tests := []struct {
		name           string
		request        *models.GetPatientRequest
		mockSetup      func(*MockPatientsService)
		expectedError  error
		expectedResult *models.PatientRequest
	}{
		{
			name:    "Success",
			request: &models.GetPatientRequest{ID: "user123"},
			mockSetup: func(m *MockPatientsService) {
				m.On("GetPatient", mock.Anything, &models.GetPatientRequest{ID: "user123"}).Return(samplePatient, nil)
			},
			expectedError:  nil,
			expectedResult: samplePatient,
		},
		{
			name:    "Not Found",
			request: &models.GetPatientRequest{ID: "nonexistent"},
			mockSetup: func(m *MockPatientsService) {
				m.On("GetPatient", mock.Anything, &models.GetPatientRequest{ID: "nonexistent"}).Return(nil, errors.New("patient not found"))
			},
			expectedError:  errors.New("patient not found"),
			expectedResult: nil,
		},
		{
			name:    "Service Error",
			request: &models.GetPatientRequest{ID: "user456"},
			mockSetup: func(m *MockPatientsService) {
				m.On("GetPatient", mock.Anything, &models.GetPatientRequest{ID: "user456"}).Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockPatientsService)
			tt.mockSetup(mockService)

			handler := &Patients{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.Get(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestPatients_GetAll(t *testing.T) {
	// Sample patients for tests
	samplePatients := []*models.PatientRequest{
		{
			ClientID:  "client123",
			FirstName: "John",
			LastName:  "Doe",
			DocType:   "DNI",
			DocNumber: "12345678",
			Gender:    "M", // Agregando género
		},
		{
			ClientID:  "client123",
			FirstName: "Jane",
			LastName:  "Smith",
			DocType:   "DNI",
			DocNumber: "87654321",
			Gender:    "F", // Agregando género
		},
	}

	// Test scenarios
	tests := []struct {
		name           string
		clientID       string
		mockSetup      func(*MockPatientsService)
		expectedError  error
		expectedResult []*models.PatientRequest
	}{
		{
			name:     "Success",
			clientID: "client123",
			mockSetup: func(m *MockPatientsService) {
				m.On("GetAllPatients", mock.Anything, "client123").Return(samplePatients, nil)
			},
			expectedError:  nil,
			expectedResult: samplePatients,
		},
		{
			name:     "Not Found",
			clientID: "nonexistent",
			mockSetup: func(m *MockPatientsService) {
				m.On("GetAllPatients", mock.Anything, "nonexistent").Return(nil, errors.New("patients not found"))
			},
			expectedError:  errors.New("patients not found"),
			expectedResult: nil,
		},
		{
			name:     "Service Error",
			clientID: "client456",
			mockSetup: func(m *MockPatientsService) {
				m.On("GetAllPatients", mock.Anything, "client456").Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockPatientsService)
			tt.mockSetup(mockService)

			handler := &Patients{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.GetAll(context.Background(), tt.clientID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestPatients_Update(t *testing.T) {
	// Test scenarios
	tests := []struct {
		name          string
		patient       *models.PatientRequest
		mockSetup     func(*MockPatientsService, *models.PatientRequest)
		expectedError error
	}{
		{
			name: "Success",
			patient: &models.PatientRequest{
				ClientID:  "client123",
				FirstName: "John",
				LastName:  "Doe",
				DocType:   "DNI",
				DocNumber: "12345678",
				Gender:    "M", // Agregando género
			},
			mockSetup: func(m *MockPatientsService, p *models.PatientRequest) {
				m.On("UpdatePatient", mock.Anything, p).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Service Error",
			patient: &models.PatientRequest{
				ClientID:  "client456",
				FirstName: "Jane",
				LastName:  "Smith",
				DocType:   "DNI",
				DocNumber: "87654321",
				Gender:    "F", // Agregando género
			},
			mockSetup: func(m *MockPatientsService, p *models.PatientRequest) {
				m.On("UpdatePatient", mock.Anything, p).Return(errors.New("service error"))
			},
			expectedError: errors.New("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockPatientsService)
			tt.mockSetup(mockService, tt.patient)

			handler := &Patients{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.Update(context.Background(), tt.patient)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.patient, result)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestPatients_Delete(t *testing.T) {
	// Test scenarios
	tests := []struct {
		name          string
		userID        string
		mockSetup     func(*MockPatientsService)
		expectedError error
	}{
		{
			name:   "Success",
			userID: "user123",
			mockSetup: func(m *MockPatientsService) {
				m.On("DeletePatient", mock.Anything, "user123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Not Found",
			userID: "nonexistent",
			mockSetup: func(m *MockPatientsService) {
				m.On("DeletePatient", mock.Anything, "nonexistent").Return(errors.New("patient not found"))
			},
			expectedError: errors.New("patient not found"),
		},
		{
			name:   "Service Error",
			userID: "user456",
			mockSetup: func(m *MockPatientsService) {
				m.On("DeletePatient", mock.Anything, "user456").Return(errors.New("service error"))
			},
			expectedError: errors.New("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockPatientsService)
			tt.mockSetup(mockService)

			handler := &Patients{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			err := handler.Delete(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockService.AssertExpectations(t)
		})
	}
}
