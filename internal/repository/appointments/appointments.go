package repository

import (
	"context"

	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

type AppointmentsRepository interface {
	Save(ctx context.Context, a *models.Appointment) error
	GetByID(ctx context.Context, id string) (*models.Appointment, error)
	GetByClientID(ctx context.Context, clientID string) ([]*models.Appointment, error)
	GetByPatientID(ctx context.Context, patientID string) ([]*models.Appointment, error)
	GetByDoctorID(ctx context.Context, doctorID string) ([]*models.Appointment, error)
	Delete(ctx context.Context, id string) error
}

type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type DynamoAppointmentsRepository struct {
	Client         DynamoDBClient
	Logger         *zap.SugaredLogger
	TableName      string
	ClientIDIndex  string
	PatientIDIndex string
	DoctorIDIndex  string
}

func New(client DynamoDBClient, logger *zap.SugaredLogger, tableName, clientIDIndex, patientIDIndex, doctorIDIndex string) AppointmentsRepository {
	return &DynamoAppointmentsRepository{
		Client:         client,
		Logger:         logger,
		TableName:      tableName,
		ClientIDIndex:  clientIDIndex,
		PatientIDIndex: patientIDIndex,
		DoctorIDIndex:  doctorIDIndex,
	}
}

func (d *DynamoAppointmentsRepository) Save(ctx context.Context, a *models.Appointment) error {
	item, err := attributevalue.MarshalMap(a)
	if err != nil {
		d.Logger.Errorw("error marshalling appointment", "error", err)
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item:      item,
	})
	return err
}

func (d *DynamoAppointmentsRepository) GetByID(ctx context.Context, id string) (*models.Appointment, error) {
	key, _ := attributevalue.MarshalMap(map[string]string{"id": id})
	resp, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key:       key,
	})
	if err != nil || resp.Item == nil {
		return nil, err
	}
	var appointment models.Appointment
	if err := attributevalue.UnmarshalMap(resp.Item, &appointment); err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (d *DynamoAppointmentsRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Appointment, error) {
	keyCond := expression.Key("client_id").Equal(expression.Value(clientID))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()

	resp, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &d.TableName,
		IndexName:                 &d.ClientIDIndex,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var results []*models.Appointment
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *DynamoAppointmentsRepository) GetByPatientID(ctx context.Context, patientID string) ([]*models.Appointment, error) {
	keyCond := expression.Key("patient_id").Equal(expression.Value(patientID))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()

	resp, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &d.TableName,
		IndexName:                 &d.PatientIDIndex,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var results []*models.Appointment
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *DynamoAppointmentsRepository) GetByDoctorID(ctx context.Context, doctorID string) ([]*models.Appointment, error) {
	keyCond := expression.Key("doctor_id").Equal(expression.Value(doctorID))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()

	resp, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &d.TableName,
		IndexName:                 &d.DoctorIDIndex,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var results []*models.Appointment
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *DynamoAppointmentsRepository) Delete(ctx context.Context, id string) error {
	key, _ := attributevalue.MarshalMap(map[string]string{"id": id})
	_, err := d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &d.TableName,
		Key:       key,
	})
	return err
}
