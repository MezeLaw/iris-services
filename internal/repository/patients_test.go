package repository

import (
	"context"
	"errors"
	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

// Mock for DynamoDB Client
type MockDynamoDBClient struct {
	mock.Mock
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

// Utility function to create a test logger
func createTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}

// TestSave tests the Save method of DynamoPatientsRepository
func TestSave(t *testing.T) {
	testCases := []struct {
		name          string
		patient       *models.Patient
		mockResponse  *dynamodb.PutItemOutput
		mockError     error
		expectedError error
	}{
		{
			name: "Success",
			patient: &models.Patient{
				ID:        "123",
				ClientID:  "client1",
				FirstName: "John",
				LastName:  "Doe",
				DocType:   "DNI",
				DocNumber: "12345678",
			},
			mockResponse:  &dynamodb.PutItemOutput{},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "DynamoDB Error",
			patient: &models.Patient{
				ID:        "123",
				ClientID:  "client1",
				FirstName: "John",
				LastName:  "Doe",
				DocType:   "DNI",
				DocNumber: "12345678",
			},
			mockResponse:  &dynamodb.PutItemOutput{},
			mockError:     errors.New("dynamodb error"),
			expectedError: errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := New(mockClient, createTestLogger(), "patients", "client_id-index", "doc_key-index")

			// Expectations
			mockClient.On("PutItem", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			err := repo.Save(context.Background(), tc.patient)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				// Verify DocKey was set correctly
				assert.Equal(t, tc.patient.DocType+"#"+tc.patient.DocNumber, tc.patient.DocKey)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestGetByID tests the GetByID method of DynamoPatientsRepository
func TestGetByID(t *testing.T) {
	// Sample patient data
	samplePatient := &models.Patient{
		ID:        "123",
		ClientID:  "client1",
		FirstName: "John",
		LastName:  "Doe",
		DocType:   "DNI",
		DocNumber: "12345678",
	}

	// Marshal the sample patient for mock response
	patientItem, _ := attributevalue.MarshalMap(samplePatient)

	testCases := []struct {
		name            string
		id              string
		mockResponse    *dynamodb.GetItemOutput
		mockError       error
		expectedPatient *models.Patient
		expectedError   error
	}{
		{
			name: "Success",
			id:   "123",
			mockResponse: &dynamodb.GetItemOutput{
				Item: patientItem,
			},
			mockError:       nil,
			expectedPatient: samplePatient,
			expectedError:   nil,
		},
		{
			name:            "Patient Not Found",
			id:              "456",
			mockResponse:    &dynamodb.GetItemOutput{Item: nil},
			mockError:       nil,
			expectedPatient: nil,
			expectedError:   nil,
		},
		{
			name:            "DynamoDB Error",
			id:              "123",
			mockResponse:    &dynamodb.GetItemOutput{},
			mockError:       errors.New("dynamodb error"),
			expectedPatient: nil,
			expectedError:   errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("GetItem", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			patient, err := repo.GetByID(context.Background(), tc.id)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, patient)
			} else if tc.expectedPatient == nil {
				assert.Nil(t, patient)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPatient.ID, patient.ID)
				assert.Equal(t, tc.expectedPatient.FirstName, patient.FirstName)
				assert.Equal(t, tc.expectedPatient.LastName, patient.LastName)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestGetByClientID tests the GetByClientID method of DynamoPatientsRepository
func TestGetByClientID(t *testing.T) {
	// Sample patients data
	samplePatients := []*models.Patient{
		{
			ID:        "123",
			ClientID:  "client1",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			ID:        "456",
			ClientID:  "client1",
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}

	// Marshal the sample patients for mock response
	var patientItems []map[string]types.AttributeValue
	for _, p := range samplePatients {
		item, _ := attributevalue.MarshalMap(p)
		patientItems = append(patientItems, item)
	}

	testCases := []struct {
		name             string
		clientID         string
		mockResponse     *dynamodb.QueryOutput
		mockError        error
		expectedPatients []*models.Patient
		expectedError    error
	}{
		{
			name:     "Success Multiple Patients",
			clientID: "client1",
			mockResponse: &dynamodb.QueryOutput{
				Items: patientItems,
			},
			mockError:        nil,
			expectedPatients: samplePatients,
			expectedError:    nil,
		},
		{
			name:     "No Patients Found",
			clientID: "client999",
			mockResponse: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{},
			},
			mockError:        nil,
			expectedPatients: []*models.Patient{},
			expectedError:    nil,
		},
		{
			name:             "DynamoDB Error",
			clientID:         "client1",
			mockResponse:     &dynamodb.QueryOutput{},
			mockError:        errors.New("dynamodb error"),
			expectedPatients: nil,
			expectedError:    errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("Query", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			patients, err := repo.GetByClientID(context.Background(), tc.clientID)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, patients)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedPatients), len(patients))
				if len(tc.expectedPatients) > 0 {
					// Check first patient details
					assert.Equal(t, tc.expectedPatients[0].ID, patients[0].ID)
					assert.Equal(t, tc.expectedPatients[0].FirstName, patients[0].FirstName)
					assert.Equal(t, tc.expectedPatients[0].LastName, patients[0].LastName)
				}
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestGetByDocument tests the GetByDocument method of DynamoPatientsRepository
func TestGetByDocument(t *testing.T) {
	// Sample patient data
	samplePatient := &models.Patient{
		ID:        "123",
		ClientID:  "client1",
		FirstName: "John",
		LastName:  "Doe",
		DocType:   "DNI",
		DocNumber: "12345678",
		DocKey:    "DNI#12345678",
	}

	// Marshal the sample patient for mock response
	patientItem, _ := attributevalue.MarshalMap(samplePatient)

	testCases := []struct {
		name            string
		docType         string
		docNumber       string
		mockResponse    *dynamodb.QueryOutput
		mockError       error
		expectedPatient *models.Patient
		expectedError   error
	}{
		{
			name:      "Success",
			docType:   "DNI",
			docNumber: "12345678",
			mockResponse: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{patientItem},
			},
			mockError:       nil,
			expectedPatient: samplePatient,
			expectedError:   nil,
		},
		{
			name:      "Patient Not Found",
			docType:   "DNI",
			docNumber: "99999999",
			mockResponse: &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{},
			},
			mockError:       nil,
			expectedPatient: nil,
			expectedError:   nil,
		},
		{
			name:            "DynamoDB Error",
			docType:         "DNI",
			docNumber:       "12345678",
			mockResponse:    &dynamodb.QueryOutput{},
			mockError:       errors.New("dynamodb error"),
			expectedPatient: nil,
			expectedError:   errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("Query", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			patient, err := repo.GetByDocument(context.Background(), tc.docType, tc.docNumber)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, patient)
			} else if tc.expectedPatient == nil {
				assert.Nil(t, patient)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPatient.ID, patient.ID)
				assert.Equal(t, tc.expectedPatient.DocType, patient.DocType)
				assert.Equal(t, tc.expectedPatient.DocNumber, patient.DocNumber)
				assert.Equal(t, tc.expectedPatient.DocKey, patient.DocKey)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestDelete tests the Delete method of DynamoPatientsRepository
func TestDelete(t *testing.T) {
	testCases := []struct {
		name          string
		id            string
		mockResponse  *dynamodb.DeleteItemOutput
		mockError     error
		expectedError error
	}{
		{
			name:          "Success",
			id:            "123",
			mockResponse:  &dynamodb.DeleteItemOutput{},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "DynamoDB Error",
			id:            "123",
			mockResponse:  &dynamodb.DeleteItemOutput{},
			mockError:     errors.New("dynamodb error"),
			expectedError: errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("DeleteItem", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			err := repo.Delete(context.Background(), tc.id)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestUpdate tests the Update method of DynamoPatientsRepository
func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		patient       *models.Patient
		mockResponse  *dynamodb.PutItemOutput
		mockError     error
		expectedError error
	}{
		{
			name: "Success",
			patient: &models.Patient{
				ID:        "123",
				ClientID:  "client1",
				FirstName: "John Updated",
				LastName:  "Doe Updated",
				DocType:   "DNI",
				DocNumber: "12345678",
			},
			mockResponse:  &dynamodb.PutItemOutput{},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "DynamoDB Error",
			patient: &models.Patient{
				ID:        "123",
				ClientID:  "client1",
				FirstName: "John",
				LastName:  "Doe",
				DocType:   "DNI",
				DocNumber: "12345678",
			},
			mockResponse:  &dynamodb.PutItemOutput{},
			mockError:     errors.New("dynamodb error"),
			expectedError: errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("PutItem", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			err := repo.Update(context.Background(), tc.patient)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				// Verify DocKey was set correctly
				assert.Equal(t, tc.patient.DocType+"#"+tc.patient.DocNumber, tc.patient.DocKey)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}

// TestGet tests the Get method of DynamoPatientsRepository
func TestGet(t *testing.T) {
	// Sample patient data
	samplePatient := &models.Patient{
		ID:        "123",
		ClientID:  "client1",
		FirstName: "John",
		LastName:  "Doe",
		DocType:   "DNI",
		DocNumber: "12345678",
	}

	// Marshal the sample patient for mock response
	patientItem, _ := attributevalue.MarshalMap(samplePatient)

	testCases := []struct {
		name            string
		id              string
		mockResponse    *dynamodb.GetItemOutput
		mockError       error
		expectedPatient *models.Patient
		expectedError   error
	}{
		{
			name: "Success",
			id:   "123",
			mockResponse: &dynamodb.GetItemOutput{
				Item: patientItem,
			},
			mockError:       nil,
			expectedPatient: samplePatient,
			expectedError:   nil,
		},
		{
			name:            "Patient Not Found",
			id:              "456",
			mockResponse:    &dynamodb.GetItemOutput{Item: nil},
			mockError:       nil,
			expectedPatient: nil,
			expectedError:   nil,
		},
		{
			name:            "DynamoDB Error",
			id:              "123",
			mockResponse:    &dynamodb.GetItemOutput{},
			mockError:       errors.New("dynamodb error"),
			expectedPatient: nil,
			expectedError:   errors.New("dynamodb error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockClient := new(MockDynamoDBClient)
			repo := DynamoPatientsRepository{
				Client:        mockClient,
				Logger:        createTestLogger(),
				TableName:     "patients",
				ClientIDIndex: "client_id-index",
				DocKeyIndex:   "doc_key-index",
			}

			// Expectations
			mockClient.On("GetItem", mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)

			// Execute
			patient, err := repo.Get(context.Background(), tc.id)

			// Assertions
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
				assert.Nil(t, patient)
			} else if tc.expectedPatient == nil {
				assert.Nil(t, patient)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPatient.ID, patient.ID)
				assert.Equal(t, tc.expectedPatient.FirstName, patient.FirstName)
				assert.Equal(t, tc.expectedPatient.LastName, patient.LastName)
			}

			// Verify mocks
			mockClient.AssertExpectations(t)
		})
	}
}
