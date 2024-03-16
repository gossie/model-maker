package domain

import (
	"context"
)

type ModelRepository interface {
	FindById(context.Context, string) (Model, error)
	FindAllByUser(context.Context, string) ([]Model, error)
	SaveModel(context.Context, string, ModelCreationRequest) (int, error)
}

type ParameterRepository interface {
	FindAllByModelId(context.Context, string) ([]Parameter, error)
	SaveParameter(context.Context, string, ParameterCreationRequest) (int, error)
	DeleteParameter(context.Context, string, string) error
	FindAllTranslations(context.Context, string) ([]Translation, error)
	SaveTranslations(context.Context, string, TranslationModificationRequest) error
	SaveValues(context.Context, string, ValueModificationRequest) error
}
