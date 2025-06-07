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

// MockAppointmentsRepository es una implementación mock de AppointmentsRepository
type MockAppointmentsRepository struct {
	mock.Mock
}

func (m *MockAppointmentsRepository) Save(ctx context.Context, a *models.Appointment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockAppointmentsRepository) GetByID(ctx context.Context, id string) (*models.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointment), args.Error(1)
}

func (m *MockAppointmentsRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Appointment, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Appointment), args.Error(1)
}

func (m *MockAppointmentsRepository) GetByPatientID(ctx context.Context, patientID string) ([]*models.Appointment, error) {
	args := m.Called(ctx, patientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Appointment), args.Error(1)
}

func (m *MockAppointmentsRepository) GetByDoctorID(ctx context.Context, doctorID string) ([]*models.Appointment, error) {
	args := m.Called(ctx, doctorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Appointment), args.Error(1)
}

func (m *MockAppointmentsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Función helper para configurar el test
func setupTest() (*Appointments, *MockAppointmentsRepository) {
	mockRepo := new(MockAppointmentsRepository)
	logger, _ := zap.NewDevelopment()
	sugarLogger := logger.Sugar()
	service := &Appointments{
		Logger:                 sugarLogger,
		AppointmentsRepository: mockRepo,
	}
	return service, mockRepo
}

// Función helper para crear una cita de ejemplo
func createSampleAppointmentRequest() *models.AppointmentRequest {
	return &models.AppointmentRequest{
		ClientID:  "client123",
		PatientID: "patient123",
		DoctorID:  "doctor123",
		Date:      time.Now().Format(time.RFC3339),
		Duration:  30,
		Status:    models.AppointmentStatusScheduled,
		Notes:     "Regular checkup",
		Metadata:  map[string]interface{}{"key": "value"},
	}
}

func createSampleAppointment(id string) *models.Appointment {
	return &models.Appointment{
		ID:        id,
		ClientID:  "client123",
		PatientID: "patient123",
		DoctorID:  "doctor123",
		Date:      time.Now().Format(time.RFC3339),
		Duration:  30,
		Status:    models.AppointmentStatusScheduled,
		Notes:     "Regular checkup",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Metadata:  map[string]interface{}{"key": "value"},
	}
}

// Tests para CreateAppointment
func TestAppointments_CreateAppointment_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	req := createSampleAppointmentRequest()

	// Expect the repository to be called with any appointment object
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Appointment")).Return(nil)

	// Execute
	result, err := service.CreateAppointment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, req, result)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_CreateAppointment_Error(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	req := createSampleAppointmentRequest()
	expectedErr := errors.New("database error")

	// Expect the repository to return an error
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Appointment")).Return(expectedErr)

	// Execute
	result, err := service.CreateAppointment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests para GetAppointment
func TestAppointments_GetAppointment_ByID_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"
	req := &models.GetAppointmentRequest{ID: appointmentID}

	expectedAppointment := createSampleAppointment(appointmentID)
	mockRepo.On("GetByID", ctx, appointmentID).Return(expectedAppointment, nil)

	// Execute
	result, err := service.GetAppointment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, appointmentID, result.ID)
	assert.Equal(t, expectedAppointment.PatientID, result.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_GetAppointment_ByID_Error(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"
	req := &models.GetAppointmentRequest{ID: appointmentID}
	expectedErr := errors.New("not found")

	mockRepo.On("GetByID", ctx, appointmentID).Return(nil, expectedErr)

	// Execute
	result, err := service.GetAppointment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_GetAppointment_ByPatientID_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	patientID := "patient123"
	req := &models.GetAppointmentRequest{PatientID: patientID}

	expectedAppointment := createSampleAppointment("appointment123")
	mockRepo.On("GetByPatientID", ctx, patientID).Return([]*models.Appointment{expectedAppointment}, nil)

	// Execute
	result, err := service.GetAppointment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAppointment.ID, result.ID)
	assert.Equal(t, patientID, result.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_GetAppointment_ByDoctorID_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	doctorID := "doctor123"
	req := &models.GetAppointmentRequest{DoctorID: doctorID}

	expectedAppointment := createSampleAppointment("appointment123")
	mockRepo.On("GetByDoctorID", ctx, doctorID).Return([]*models.Appointment{expectedAppointment}, nil)

	// Execute
	result, err := service.GetAppointment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAppointment.ID, result.ID)
	assert.Equal(t, doctorID, result.DoctorID)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_GetAppointment_InvalidParams(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()
	req := &models.GetAppointmentRequest{} // Empty request

	// Execute
	result, err := service.GetAppointment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid parameters")
}

// Tests para GetAllAppointments
func TestAppointments_GetAllAppointments_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	clientID := "client123"

	appointment1 := createSampleAppointment("appointment1")
	appointment2 := createSampleAppointment("appointment2")
	appointments := []*models.Appointment{appointment1, appointment2}

	mockRepo.On("GetByClientID", ctx, clientID).Return(appointments, nil)

	// Execute
	result, err := service.GetAllAppointments(ctx, clientID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, appointment1.ID, result[0].ID)
	assert.Equal(t, appointment2.ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_GetAllAppointments_EmptyClientID(t *testing.T) {
	// Setup
	service, _ := setupTest()
	ctx := context.Background()

	// Execute
	result, err := service.GetAllAppointments(ctx, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "client-id cannot be empty")
}

func TestAppointments_GetAllAppointments_RepositoryError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	clientID := "client123"
	expectedErr := errors.New("database error")

	mockRepo.On("GetByClientID", ctx, clientID).Return(nil, expectedErr)

	// Execute
	result, err := service.GetAllAppointments(ctx, clientID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tests para UpdateAppointment
func TestAppointments_UpdateAppointment_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"

	req := createSampleAppointmentRequest()
	req.ID = appointmentID

	existingAppointment := createSampleAppointment(appointmentID)

	mockRepo.On("GetByID", ctx, appointmentID).Return(existingAppointment, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Appointment")).Return(nil)

	// Execute
	err := service.UpdateAppointment(ctx, req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_UpdateAppointment_NotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"

	req := createSampleAppointmentRequest()
	req.ID = appointmentID

	mockRepo.On("GetByID", ctx, appointmentID).Return(nil, errors.New("not found"))

	// Execute
	err := service.UpdateAppointment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find appointment")
	mockRepo.AssertExpectations(t)
}

func TestAppointments_UpdateAppointment_SaveError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"

	req := createSampleAppointmentRequest()
	req.ID = appointmentID

	existingAppointment := createSampleAppointment(appointmentID)
	expectedErr := errors.New("database error")

	mockRepo.On("GetByID", ctx, appointmentID).Return(existingAppointment, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*models.Appointment")).Return(expectedErr)

	// Execute
	err := service.UpdateAppointment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update appointment")
	mockRepo.AssertExpectations(t)
}

// Tests para DeleteAppointment
func TestAppointments_DeleteAppointment_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"

	mockRepo.On("GetByID", ctx, appointmentID).Return(createSampleAppointment(appointmentID), nil)
	mockRepo.On("Delete", ctx, appointmentID).Return(nil)

	// Execute
	err := service.DeleteAppointment(ctx, appointmentID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointments_DeleteAppointment_NotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"

	mockRepo.On("GetByID", ctx, appointmentID).Return(nil, errors.New("not found"))

	// Execute
	err := service.DeleteAppointment(ctx, appointmentID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find appointment")
	mockRepo.AssertExpectations(t)
}

func TestAppointments_DeleteAppointment_DeleteError(t *testing.T) {
	// Setup
	service, mockRepo := setupTest()
	ctx := context.Background()
	appointmentID := "appointment123"
	expectedErr := errors.New("database error")

	mockRepo.On("GetByID", ctx, appointmentID).Return(createSampleAppointment(appointmentID), nil)
	mockRepo.On("Delete", ctx, appointmentID).Return(expectedErr)

	// Execute
	err := service.DeleteAppointment(ctx, appointmentID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete appointment")
	mockRepo.AssertExpectations(t)
}

// Test para validación de status
func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  models.AppointmentStatus
		wantErr bool
	}{
		{
			name:    "Valid Status - Scheduled",
			status:  models.AppointmentStatusScheduled,
			wantErr: false,
		},
		{
			name:    "Valid Status - In Progress",
			status:  models.AppointmentStatusInProgress,
			wantErr: false,
		},
		{
			name:    "Valid Status - Completed",
			status:  models.AppointmentStatusCompleted,
			wantErr: false,
		},
		{
			name:    "Valid Status - Cancelled",
			status:  models.AppointmentStatusCancelled,
			wantErr: false,
		},
		{
			name:    "Invalid Status",
			status:  "INVALID_STATUS",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStatus(tt.status)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
