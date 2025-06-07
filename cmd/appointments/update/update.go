package main

import (
	"context"
	"encoding/json"
	"time"

	handler "github.com/MezeLaw/iris-services/internal/handler/appointments"
	"github.com/MezeLaw/iris-services/internal/models"
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
		var request models.AppointmentRequest
		if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
			sugar.Errorf("Error unmarshalling request: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"invalid request body"}`}, nil
		}

		// Asegurarse de que el ID esté presente
		if request.ID == "" {
			sugar.Error("Missing appointment ID in request")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"missing appointment ID"}`}, nil
		}

		request.UpdatedAt = time.Now().Format(time.RFC3339)

		updated, err := h.Update(ctx, &request)
		if err != nil {
			sugar.Errorf("Error updating appointment: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not update appointment"}`}, nil
		}

		respBody, _ := json.Marshal(updated)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
