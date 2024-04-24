package rest

type RenderModel struct {
	Id   int
	Name string
}

type RenderParameter struct {
	Id      int
	ModelId int
	Name    string
	Values  []string
}

type RenderConstraint struct {
}

type ModelCatalogRenderContext struct {
	Models []RenderModel
}

type ModelRenderContext struct {
	Model       RenderModel
	Parameters  []RenderParameter
	Constraints []RenderConstraint
}

// TODO: delete
type parameterCreationResponse struct {
	ParameterId int
}
