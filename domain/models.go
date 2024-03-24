package domain

import (
	configurationmodel "github.com/gossie/configuration-model"
)

type ModelCreationRequest struct {
	Name string `json:"name"`
}

type Model struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Translation string       `json:"translation"`
	Constraints []Constraint `json:"constraints"`
}

type ConstraintCreationRequest struct {
	Type          configurationmodel.ConstraintType `json:"type"`
	FromId        int                               `json:"fromId"`
	FromValueId   int                               `json:"fromValueId"`
	TargetId      int                               `json:"targetId"`
	TargetValueId int                               `json:"targetValueId"`
}

type ConstraintCreationResponse struct {
	Id int `json:"id"`
}

type Constraint struct {
	Id            int                               `json:"id"`
	Type          configurationmodel.ConstraintType `json:"type"`
	FromId        int                               `json:"fromId"`
	FromValueId   int                               `json:"fromValueId"`
	TargetId      int                               `json:"targetId"`
	TargetValueId int                               `json:"targetValueId"`
}

type ParameterCreationRequest struct {
	Name      string                       `json:"name"`
	ValueType configurationmodel.ValueType `json:"valueType"`
}

type Parameter struct {
	Id          int                          `json:"id"`
	Name        string                       `json:"name"`
	Translation string                       `json:"translation"`
	ValueType   configurationmodel.ValueType `json:"valueType"`
	Value       ParameterValue               `json:"value"`
}

type ParameterValue struct {
	Values []Value `json:"values"`
}

type Value struct {
	Id          int    `json:"id"`
	Value       string `json:"value"`
	Translation string `json:"translation"`
}

type ValueModificationRequest struct {
	NewValues     []string `json:"newValues"`
	UpdatedValues []Value  `json:"updatedValues"`
}

type TranslationModificationRequest struct {
	NewTranslations     []Translation `json:"newTranslations"`
	UpdatedTranslations []Translation `json:"updatedTranslations"`
}

type Translation struct {
	Id       int    `json:"id"`
	Field    string `json:"field"`
	Language string `json:"language"`
	Value    string `json:"value"`
}
