package logging

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction() // Para prod
	// logger, err := zap.NewDevelopment() // Para local con colores
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
