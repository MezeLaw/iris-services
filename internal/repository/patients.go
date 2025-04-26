package repository

import (
	"context"
	"fmt"
	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

type PatientsRepository interface {
	Save(ctx context.Context, p *models.Patient) error
	GetByID(ctx context.Context, id string) (*models.Patient, error)
	GetByClientID(ctx context.Context, clientID string) ([]*models.Patient, error)
	GetByDocument(ctx context.Context, docType, docNumber string) (*models.Patient, error)
	Delete(ctx context.Context, id string) error
}

type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type DynamoPatientsRepository struct {
	Client        DynamoDBClient
	Logger        *zap.SugaredLogger
	TableName     string
	ClientIDIndex string
	DocKeyIndex   string
}

func New(client DynamoDBClient, logger *zap.SugaredLogger, tableName, clientIDIndex, docKeyIndex string) PatientsRepository {
	return &DynamoPatientsRepository{
		Client:        client,
		Logger:        logger,
		TableName:     tableName,
		ClientIDIndex: clientIDIndex,
		DocKeyIndex:   docKeyIndex,
	}
}

func (d *DynamoPatientsRepository) Save(ctx context.Context, p *models.Patient) error {
	p.DocKey = fmt.Sprintf("%s#%s", p.DocType, p.DocNumber)
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		d.Logger.Errorw("error marshalling patient", "error", err)
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item:      item,
	})
	return err
}

func (d *DynamoPatientsRepository) Get(ctx context.Context, id string) (*models.Patient, error) {
	key, _ := attributevalue.MarshalMap(map[string]string{"id": id})
	resp, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key:       key,
	})
	if err != nil || resp.Item == nil {
		return nil, err
	}
	var patient models.Patient
	if err := attributevalue.UnmarshalMap(resp.Item, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (d *DynamoPatientsRepository) Update(ctx context.Context, p *models.Patient) error {
	p.DocKey = fmt.Sprintf("%s#%s", p.DocType, p.DocNumber)
	item, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.TableName,
		Item:      item,
	})
	return err
}

func (d *DynamoPatientsRepository) Delete(ctx context.Context, id string) error {
	key, _ := attributevalue.MarshalMap(map[string]string{"id": id})
	_, err := d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &d.TableName,
		Key:       key,
	})
	return err
}
func (d *DynamoPatientsRepository) GetByID(ctx context.Context, id string) (*models.Patient, error) {
	key, _ := attributevalue.MarshalMap(map[string]string{"id": id})
	resp, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.TableName,
		Key:       key,
	})
	if err != nil || resp.Item == nil {
		return nil, err
	}
	var patient models.Patient
	if err := attributevalue.UnmarshalMap(resp.Item, &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}

func (d *DynamoPatientsRepository) GetByClientID(ctx context.Context, clientID string) ([]*models.Patient, error) {
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

	var results []*models.Patient
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (d *DynamoPatientsRepository) GetByDocument(ctx context.Context, docType, docNumber string) (*models.Patient, error) {
	docKey := fmt.Sprintf("%s#%s", docType, docNumber)
	keyCond := expression.Key("doc_key").Equal(expression.Value(docKey))
	expr, _ := expression.NewBuilder().WithKeyCondition(keyCond).Build()

	resp, err := d.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &d.TableName,
		IndexName:                 &d.DocKeyIndex,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1),
	})
	if err != nil || len(resp.Items) == 0 {
		return nil, err
	}
	var patient models.Patient
	if err := attributevalue.UnmarshalMap(resp.Items[0], &patient); err != nil {
		return nil, err
	}
	return &patient, nil
}
