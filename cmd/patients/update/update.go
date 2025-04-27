package main

import (
	"context"
	"encoding/json"
	"time"

	handler "github.com/MezeLaw/iris-services/internal/handler/patients"
	"github.com/MezeLaw/iris-services/internal/models"
	repository "github.com/MezeLaw/iris-services/internal/repository/patients"
	service "github.com/MezeLaw/iris-services/internal/service/patients"
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
		var request models.PatientRequest
		if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
			sugar.Errorf("Error unmarshalling request: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"invalid request body"}`}, nil
		}

		// Ensure ID is provided for update
		if request.ID == "" {
			sugar.Error("Missing patient ID in update request")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"patient ID is required for update"}`}, nil
		}

		// Set update timestamp
		request.UpdatedAt = time.Now().Format(time.RFC3339)

		updated, err := h.Update(ctx, &request)
		if err != nil {
			sugar.Errorf("Error updating patient: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not update patient"}`}, nil
		}

		respBody, _ := json.Marshal(updated)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
