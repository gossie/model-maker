package persistence

import (
	"context"

	"github.com/gossie/modelling-service/domain"
)

type ModelRepository interface {
	FindById(ctx context.Context, id string) (domain.Model, error)
	FindAllByUser(ctx context.Context, id string) ([]domain.Model, error)
	SaveModel(ctx context.Context, userEmail string, cmr domain.ModelCreationRequest) (int, error)
}

type ParameterRepository interface {
	FindAllByModelId(ctx context.Context, modelId string) ([]domain.Parameter, error)
	SaveParameter(ctx context.Context, modelId string, pmr domain.ParameterCreationRequest) (int, error)
	DeleteParameter(ctx context.Context, modelId string, parameterId string) error
	FindAllTranslations(ctx context.Context, parameterId string) ([]domain.Translation, error)
	SaveTranslations(ctx context.Context, parameterId string, tcr domain.TranslationModificationRequest) error
}
