package server

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *server) findModel(ctx context.Context, modelId string) (model, error) {
	result := model{}

	var id int
	var name string
	row := s.db.QueryRowContext(ctx, "SELECT * FROM models WHERE id = $1", modelId)
	err := row.Scan(&id, &name)

	result.Id = id
	result.Name = name
	return result, err
}

func (s *server) findAllModelsByUser(ctx context.Context, userEmail string) ([]model, error) {
	sqlStatement := `
		SELECT id, name
		FROM models
		JOIN model_user_relations
		ON models.id = model_user_relations.modelid AND model_user_relations.userid = (SELECT id FROM users WHERE email = $1);
	`
	rows, err := s.db.QueryContext(ctx, sqlStatement, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]model, 0)

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		models = append(models, model{id, name})
	}

	return models, nil
}

func (s *server) saveModel(ctx context.Context, userEmail string, cmr modelCreationRequest) (int, error) {
	tx, _ := s.db.BeginTx(ctx, nil)

	var userId int
	row := tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", userEmail)
	err := row.Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}

	var modelId int
	err = s.db.QueryRowContext(ctx, "INSERT INTO models (name) VALUES ($1) RETURNING id", cmr.Name).Scan(&modelId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("created model with ID %v", modelId))

	_, err = s.db.ExecContext(ctx, "INSERT INTO model_user_relations (modelId, userId) VALUES ($1, $2)", modelId, userId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return modelId, nil
}

func (s *server) findAllParametersByModelId(ctx context.Context, modelId string) ([]parameter, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name FROM parameters WHERE modelId = $1", modelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parameters := make([]parameter, 0)

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, parameter{id, name})
	}

	return parameters, nil
}

func (s *server) saveParameter(ctx context.Context, modelId string, pmr parameterCreationRequest) (int, error) {
	var parameterId int
	err := s.db.QueryRowContext(ctx, "INSERT INTO parameters (name, modelId) VALUES ($1, $2) RETURNING id", pmr.Name, modelId).Scan(&parameterId)
	return parameterId, err
}
