package handler

import (
	"context"
	"github.com/MezeLaw/iris-services/internal/models"
)

type PatientsHandler interface {
	Create(context.Context, models.Patient) (string, error)
	Get(ctx context.Context)
	Update(ctx context.Context)
	Delete(ctx context.Context)
}

type Patients struct {
}

func New() PatientsHandler {
	return &Patients{}
}

func (p Patients) Create(ctx context.Context, patient models.Patient) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p Patients) Get(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p Patients) Update(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p Patients) Delete(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
