package main

import (
	"context"

	"github.com/MezeLaw/iris-services/internal/handler"
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
		if patientID == "" {
			sugar.Error("Missing patient ID in delete request")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"patient ID is required for deletion"}`}, nil
		}

		err := h.Delete(ctx, patientID)
		if err != nil {
			sugar.Errorf("Error deleting patient: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not delete patient"}`}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
