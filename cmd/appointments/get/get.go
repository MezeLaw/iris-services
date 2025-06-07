package main

import (
	"context"
	"encoding/json"

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
		// Crear el request con los parámetros disponibles
		getRequest := &models.GetAppointmentRequest{}

		// Intentar obtener el ID del path
		if id := req.PathParameters["id"]; id != "" {
			getRequest.ID = id
		}

		// Intentar obtener otros parámetros de la query string
		if clientID := req.QueryStringParameters["clientId"]; clientID != "" {
			getRequest.ClientID = clientID
		}
		if patientID := req.QueryStringParameters["patientId"]; patientID != "" {
			getRequest.PatientID = patientID
		}
		if doctorID := req.QueryStringParameters["doctorId"]; doctorID != "" {
			getRequest.DoctorID = doctorID
		}

		// Verificar que al menos un parámetro de búsqueda esté presente
		if getRequest.ID == "" && getRequest.ClientID == "" && getRequest.PatientID == "" && getRequest.DoctorID == "" {
			sugar.Error("No search parameters provided")
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: `{"error":"at least one search parameter is required"}`}, nil
		}

		appointment, err := h.Get(ctx, getRequest)
		if err != nil {
			sugar.Errorf("Error retrieving appointment: %v", err.Error())
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"could not retrieve appointment"}`}, nil
		}

		respBody, _ := json.Marshal(appointment)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(respBody),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	})
}
