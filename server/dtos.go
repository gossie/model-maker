package server

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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type parameterCreationRequest struct {
	Name string `json:"name"`
}

type parameterCreationResponse struct {
	ParameterId int `json:"parameterId"`
}

type parameter struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
