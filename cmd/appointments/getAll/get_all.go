package main

import (
	"context"
	"encoding/json"

	handler "github.com/MezeLaw/iris-services/internal/handler/appointments"
	repository "github.com/MezeLaw/iris-services/internal/repository/appointments"
	service "github.com/MezeLaw/iris-services/internal/service/appointments"
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

	repo := repository.New(dynamoClient, sugar, "AppointmentsTable", "client_id_index", "patient_id_index", "doctor_id_index")
	svc := service.New(sugar, repo)
	h := handler.New(svc, sugar)

	lambda.Start(func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		clientID := req.QueryStringParameters["clientId"]
		if clientID == "" {
			sugar.Error("Missing clientId parameter in request")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"missing clientId parameter"}`}, nil
		}

		appointments, err := h.GetAll(ctx, clientID)
		if err != nil {
			sugar.Errorf("Error retrieving appointments: %v", err)
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not retrieve appointments"}`}, nil
		}

		respBody, _ := json.Marshal(appointments)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
