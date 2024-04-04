package domain

import (
	"context"
)

type ModelRepository interface {
	FindById(context.Context, int) (Model, error)
	FindAllByUser(context.Context, string) ([]Model, error)
	SaveModel(context.Context, string, ModelCreationRequest) (int, error)
}

type ParameterRepository interface {
	FindAllByModelId(context.Context, int) ([]Parameter, error)
	SaveParameter(context.Context, int, ParameterCreationRequest) (int, error)
	DeleteParameter(context.Context, int, int) error
	FindAllTranslations(context.Context, string) ([]Translation, error)
	SaveTranslations(context.Context, string, TranslationModificationRequest) error
	SaveValues(context.Context, string, ValueModificationRequest) error
}

type ConstraintRepository interface {
	SaveConstraint(context.Context, string, ConstraintCreationRequest) (int, error)
	DeleteConstraint(context.Context, string, string) error
}
