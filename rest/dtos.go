package rest

type loginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type modelCreationResponse struct {
	ModelId int `json:"modelId"`
}

type parameterCreationResponse struct {
	ParameterId int `json:"parameterId"`
}
