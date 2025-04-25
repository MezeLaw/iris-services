package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type DynamoPatientsRepository struct {
	Client    DynamoDBClient
	TableName string
}

func New(client DynamoDBClient, tableName string) DynamoDBClient {
	return &DynamoPatientsRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (d DynamoPatientsRepository) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d DynamoPatientsRepository) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (d DynamoPatientsRepository) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	//TODO implement me
	panic("implement me")
}
