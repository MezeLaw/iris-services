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
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		sugar.Fatalf("error loading AWS config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)

	repo := repository.New(dynamoClient, sugar, "AppointmentsTable", "client_id_index", "doc_key_index")
	svc := service.New(sugar, repo)
	h := handler.New(svc, sugar)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var request models.AppointmentRequest
		if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
			sugar.Errorf("Error unmarshalling request: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"invalid request body"}`}, nil
		}

		now := time.Now().Format(time.RFC3339)
		request.CreatedAt = now
		request.UpdatedAt = now

		created, err := h.Create(ctx, &request)
		if err != nil {
			sugar.Errorf("Error creating appointment: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not create appointment"}`}, nil
		}

		respBody, _ := json.Marshal(created)
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
