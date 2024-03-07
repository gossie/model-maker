package server

import "github.com/gossie/configurator"

type loginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type modelCreationRequest struct {
	Name string `json:"name"`
}

type modelCreationResponse struct {
	ModelId int `json:"modelId"`
}

type model struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type parameterCreationRequest struct {
	Name      string                 `json:"name"`
	ValueType configurator.ValueType `json:"valueType"`
}

type parameterCreationResponse struct {
	ParameterId int `json:"parameterId"`
}

type parameter struct {
	Id          int                    `json:"id"`
	Name        string                 `json:"name"`
	Translation string                 `json:"translation"`
	ValueType   configurator.ValueType `json:"valueType"`
	Value       value                  `json:"value"`
}

type value struct {
	Values []string `json:"values"`
}
