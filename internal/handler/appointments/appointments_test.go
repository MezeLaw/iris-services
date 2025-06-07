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

// MockAppointmentsService implementa la interfaz AppointmentsService para los tests
type MockAppointmentsService struct {
	mock.Mock
}

func (m *MockAppointmentsService) CreateAppointment(ctx context.Context, appointment *models.AppointmentRequest) (*models.AppointmentRequest, error) {
	args := m.Called(ctx, appointment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AppointmentRequest), args.Error(1)
}

func (m *MockAppointmentsService) GetAppointment(ctx context.Context, req *models.GetAppointmentRequest) (*models.AppointmentRequest, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AppointmentRequest), args.Error(1)
}

func (m *MockAppointmentsService) GetAllAppointments(ctx context.Context, clientID string) ([]*models.AppointmentRequest, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AppointmentRequest), args.Error(1)
}

func (m *MockAppointmentsService) UpdateAppointment(ctx context.Context, appointment *models.AppointmentRequest) error {
	args := m.Called(ctx, appointment)
	return args.Error(0)
}

func (m *MockAppointmentsService) DeleteAppointment(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	// Arrange
	logger := zaptest.NewLogger(t).Sugar()
	mockService := new(MockAppointmentsService)

	// Act
	handler := New(mockService, logger)

	// Assert
	assert.NotNil(t, handler)
	assert.IsType(t, &Appointments{}, handler)
}

func TestAppointments_Create(t *testing.T) {
	// Sample appointment for tests
	sampleAppointment := &models.AppointmentRequest{
		ClientID:  "client123",
		PatientID: "patient123",
		DoctorID:  "doctor123",
		Status:    models.AppointmentStatusScheduled,
	}

	// Error appointment with invalid status
	errorAppointment := &models.AppointmentRequest{
		ClientID:  "client456",
		PatientID: "patient456",
		DoctorID:  "doctor456",
		Status:    "INVALID_STATUS",
	}

	tests := []struct {
		name           string
		appointment    *models.AppointmentRequest
		mockSetup      func(*MockAppointmentsService, *models.AppointmentRequest)
		expectedError  error
		expectedResult *models.AppointmentRequest
	}{
		{
			name:        "Success",
			appointment: sampleAppointment,
			mockSetup: func(m *MockAppointmentsService, a *models.AppointmentRequest) {
				m.On("CreateAppointment", mock.Anything, a).Return(a, nil)
			},
			expectedError:  nil,
			expectedResult: sampleAppointment,
		},
		{
			name:        "Service Error",
			appointment: sampleAppointment,
			mockSetup: func(m *MockAppointmentsService, a *models.AppointmentRequest) {
				m.On("CreateAppointment", mock.Anything, a).Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
		{
			name:        "Invalid Status",
			appointment: errorAppointment,
			mockSetup:   func(m *MockAppointmentsService, a *models.AppointmentRequest) {},
			expectedError: errors.New("invalid status value: INVALID_STATUS. Must be one of: SCHEDULED, IN_PROGRESS, COMPLETED, " +
				"CANCELLED"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockAppointmentsService)
			tt.mockSetup(mockService, tt.appointment)

			handler := &Appointments{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.Create(context.Background(), tt.appointment)

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

func TestAppointments_Get(t *testing.T) {
	// Sample appointment for tests
	sampleAppointment := &models.AppointmentRequest{
		ClientID:  "client123",
		PatientID: "patient123",
		DoctorID:  "doctor123",
		Status:    models.AppointmentStatusScheduled,
	}

	// Test scenarios
	tests := []struct {
		name           string
		request        *models.GetAppointmentRequest
		mockSetup      func(*MockAppointmentsService)
		expectedError  error
		expectedResult *models.AppointmentRequest
	}{
		{
			name:    "Success",
			request: &models.GetAppointmentRequest{ID: "appointment123"},
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAppointment", mock.Anything, &models.GetAppointmentRequest{ID: "appointment123"}).Return(sampleAppointment, nil)
			},
			expectedError:  nil,
			expectedResult: sampleAppointment,
		},
		{
			name:    "Not Found",
			request: &models.GetAppointmentRequest{ID: "nonexistent"},
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAppointment", mock.Anything, &models.GetAppointmentRequest{ID: "nonexistent"}).Return(nil, errors.New("appointment not found"))
			},
			expectedError:  errors.New("appointment not found"),
			expectedResult: nil,
		},
		{
			name:    "Service Error",
			request: &models.GetAppointmentRequest{ID: "appointment456"},
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAppointment", mock.Anything, &models.GetAppointmentRequest{ID: "appointment456"}).Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockAppointmentsService)
			tt.mockSetup(mockService)

			handler := &Appointments{
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

func TestAppointments_GetAll(t *testing.T) {
	// Sample appointments for tests
	sampleAppointments := []*models.AppointmentRequest{
		{
			ClientID:  "client123",
			PatientID: "patient123",
			DoctorID:  "doctor123",
			Status:    models.AppointmentStatusScheduled,
		},
		{
			ClientID:  "client123",
			PatientID: "patient456",
			DoctorID:  "doctor456",
			Status:    models.AppointmentStatusCompleted,
		},
	}

	// Test scenarios
	tests := []struct {
		name           string
		clientID       string
		mockSetup      func(*MockAppointmentsService)
		expectedError  error
		expectedResult []*models.AppointmentRequest
	}{
		{
			name:     "Success",
			clientID: "client123",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAllAppointments", mock.Anything, "client123").Return(sampleAppointments, nil)
			},
			expectedError:  nil,
			expectedResult: sampleAppointments,
		},
		{
			name:     "Not Found",
			clientID: "nonexistent",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAllAppointments", mock.Anything, "nonexistent").Return(nil, errors.New("appointments not found"))
			},
			expectedError:  errors.New("appointments not found"),
			expectedResult: nil,
		},
		{
			name:     "Service Error",
			clientID: "client456",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("GetAllAppointments", mock.Anything, "client456").Return(nil, errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockAppointmentsService)
			tt.mockSetup(mockService)

			handler := &Appointments{
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

func TestAppointments_Update(t *testing.T) {
	// Sample appointment for tests
	sampleAppointment := &models.AppointmentRequest{
		ID:        "appointment123",
		ClientID:  "client123",
		PatientID: "patient123",
		DoctorID:  "doctor123",
		Status:    models.AppointmentStatusScheduled,
	}

	// Test scenarios
	tests := []struct {
		name           string
		appointment    *models.AppointmentRequest
		mockSetup      func(*MockAppointmentsService)
		expectedError  error
		expectedResult *models.AppointmentRequest
	}{
		{
			name:        "Success",
			appointment: sampleAppointment,
			mockSetup: func(m *MockAppointmentsService) {
				m.On("UpdateAppointment", mock.Anything, sampleAppointment).Return(nil)
			},
			expectedError:  nil,
			expectedResult: sampleAppointment,
		},
		{
			name: "Invalid Status",
			appointment: &models.AppointmentRequest{
				ID:        "appointment123",
				ClientID:  "client123",
				PatientID: "patient123",
				DoctorID:  "doctor123",
				Status:    "INVALID_STATUS",
			},
			mockSetup:      func(m *MockAppointmentsService) {},
			expectedError:  errors.New("invalid status value: INVALID_STATUS. Must be one of: SCHEDULED, IN_PROGRESS, COMPLETED, CANCELLED"),
			expectedResult: nil,
		},
		{
			name:        "Service Error",
			appointment: sampleAppointment,
			mockSetup: func(m *MockAppointmentsService) {
				m.On("UpdateAppointment", mock.Anything, sampleAppointment).Return(errors.New("service error"))
			},
			expectedError:  errors.New("service error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockAppointmentsService)
			tt.mockSetup(mockService)

			handler := &Appointments{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			result, err := handler.Update(context.Background(), tt.appointment)

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

func TestAppointments_Delete(t *testing.T) {
	// Test scenarios
	tests := []struct {
		name          string
		appointmentID string
		mockSetup     func(*MockAppointmentsService)
		expectedError error
	}{
		{
			name:          "Success",
			appointmentID: "appointment123",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("DeleteAppointment", mock.Anything, "appointment123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "Not Found",
			appointmentID: "nonexistent",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("DeleteAppointment", mock.Anything, "nonexistent").Return(errors.New("appointment not found"))
			},
			expectedError: errors.New("appointment not found"),
		},
		{
			name:          "Service Error",
			appointmentID: "appointment456",
			mockSetup: func(m *MockAppointmentsService) {
				m.On("DeleteAppointment", mock.Anything, "appointment456").Return(errors.New("service error"))
			},
			expectedError: errors.New("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := zaptest.NewLogger(t).Sugar()
			mockService := new(MockAppointmentsService)
			tt.mockSetup(mockService)

			handler := &Appointments{
				Service: mockService,
				Logger:  logger,
			}

			// Act
			err := handler.Delete(context.Background(), tt.appointmentID)

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
