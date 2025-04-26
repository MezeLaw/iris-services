package main

import (
	"context"
	"encoding/json"

	"github.com/MezeLaw/iris-services/internal/handler"
	"github.com/MezeLaw/iris-services/internal/models"
	"github.com/MezeLaw/iris-services/internal/repository"
	"github.com/MezeLaw/iris-services/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		sugar.Fatalf("error loading AWS config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	repo := repository.New(dynamoClient, sugar, "PatientsTable", "client_id_index", "doc_key_index")
	svc := service.New(sugar, repo)
	h := handler.New(svc, sugar)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		patientID := req.QueryStringParameters["id"]
		docType := req.QueryStringParameters["docType"]
		docNumber := req.QueryStringParameters["docNumber"]

		if patientID == "" && (docType == "" || docNumber == "") {
			sugar.Errorf("Missing required parameters in request")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"missing required parameters"}`}, nil
		}

		request := &models.GetPatientRequest{
			ID:        patientID,
			DocType:   docType,
			DocNumber: docNumber,
		}

		patient, err := h.Get(ctx, request)
		if err != nil {
			sugar.Errorf("Error retrieving patient: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not retrieve patient"}`}, nil
		}

		respBody, _ := json.Marshal(patient)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
