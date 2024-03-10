package domain

import "github.com/gossie/configurator"

type ModelCreationRequest struct {
	Name string `json:"name"`
}

type Model struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type ParameterCreationRequest struct {
	Name      string                 `json:"name"`
	ValueType configurator.ValueType `json:"valueType"`
}

type Parameter struct {
	Id          int                    `json:"id"`
	Name        string                 `json:"name"`
	Translation string                 `json:"translation"`
	ValueType   configurator.ValueType `json:"valueType"`
	Value       Value                  `json:"value"`
}

type Value struct {
	Values []string `json:"values"`
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
