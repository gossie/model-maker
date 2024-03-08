package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gossie/configurator"
	"github.com/gossie/modelling-service/middleware"
)

func (s *server) findModel(ctx context.Context, modelId string) (model, error) {
	sqlStatement := `
		SELECT m.id, m.name, t.translation FROM models m
		JOIN model_translations t
		ON m.id = t.modelId
		WHERE m.id = $1 AND t.language = $2
	`
	var id int
	var name string
	var translation sql.NullString
	row := s.db.QueryRowContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey))
	err := row.Scan(&id, &name, &translation)

	return model{id, name, translation.String}, err
}

func (s *server) findAllModelsByUser(ctx context.Context, userEmail string) ([]model, error) {
	sqlStatement := `
		SELECT m.id, m.name, t.translation
		FROM models m
		JOIN model_user_relations mur
		ON m.id = mur.modelid
		JOIN model_translations t ON m.id = t.modelId
		WHERE mur.userid = (SELECT id FROM users WHERE email = $1) AND t.language = $2;
	`
	rows, err := s.db.QueryContext(ctx, sqlStatement, userEmail, ctx.Value(middleware.LanguageKey))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]model, 0)

	for rows.Next() {
		var id int
		var name string
		var translation sql.NullString
		err = rows.Scan(&id, &name, &translation)
		if err != nil {
			return nil, err
		}
		models = append(models, model{id, name, translation.String})
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
	sqlStatement := `
		SELECT p.id, p.name, p.valueType, p.value, t.translation
		FROM parameters p
		LEFT JOIN parameter_translations t
		ON p.id = t.parameterId
		WHERE p.modelId = $1 AND (t.language = $2 OR t.language IS NULL)
	`

	rows, err := s.db.QueryContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parameters := make([]parameter, 0)

	for rows.Next() {
		var id int
		var name string
		var valueType configurator.ValueType
		var paramValue sql.NullString
		var translation sql.NullString
		err = rows.Scan(&id, &name, &valueType, &paramValue, &translation)
		if err != nil {
			return nil, err
		}

		stringValues := make([]string, 0)
		if paramValue.Valid {
			decoder := json.NewDecoder(strings.NewReader(paramValue.String))
			err = nil
			switch valueType {
			case configurator.IntSetType:
				values := make([]int, 0)
				err = decoder.Decode(&values)
				for _, v := range values {
					stringValues = append(stringValues, fmt.Sprint(v))
				}
			case configurator.StringSetType:
				err = decoder.Decode(&stringValues)
			}

			if err != nil {
				return nil, err
			}
		}

		parameters = append(parameters, parameter{id, name, translation.String, valueType, value{stringValues}})
	}

	return parameters, nil
}

func (s *server) saveParameter(ctx context.Context, modelId string, pmr parameterCreationRequest) (int, error) {
	var parameterId int
	err := s.db.QueryRowContext(ctx, "INSERT INTO parameters (name, valueType, modelId) VALUES ($1, $2, $3) RETURNING id", pmr.Name, pmr.ValueType, modelId).Scan(&parameterId)
	return parameterId, err
}

func (s *server) deleteParameterFromModel(ctx context.Context, modelId string, parameterId string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM parameters WHERE id = $1 AND modelId = $2", parameterId, modelId)
	return err
}

func (s *server) findAllParameterTranslations(ctx context.Context, parameterId string) ([]translation, error) {
	sqlStatement := `
		SELECT id, field, language, translation
		FROM parameter_translations
		WHERE parameterId = $1
	`

	rows, err := s.db.QueryContext(ctx, sqlStatement, parameterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	translations := make([]translation, 0)

	for rows.Next() {
		var id int
		var field string
		var language string
		var value string
		err = rows.Scan(&id, &field, &language, &value)
		if err != nil {
			return nil, err
		}

		translations = append(translations, translation{id, field, language, value})
	}

	return translations, nil
}

func (s *server) saveParameterTranslations(ctx context.Context, parameterId string, tcr translationModificationRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if len(tcr.NewTranslations) > 0 {
		args := make([]any, 0, len(tcr.NewTranslations)*4)
		valueStrings := make([]string, 0, len(tcr.NewTranslations))
		for _, translation := range tcr.NewTranslations {
			valueStrings = append(valueStrings, fmt.Sprintf("($%v, $%v, $%v, $%v)", len(args)+1, len(args)+2, len(args)+3, len(args)+4))
			args = append(args, parameterId, translation.Field, translation.Language, translation.Value)
		}

		sqlStatement := `
		INSERT INTO parameter_translations (parameterId, field, language, translation)
		VALUES ` + strings.Join(valueStrings, ", ")
		_, err = tx.ExecContext(ctx, sqlStatement, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(tcr.UpdatedTranslations) > 0 {
		for _, translation := range tcr.UpdatedTranslations {
			sqlStatement := `
				UPDATE parameter_translations
				SET language = $1, translation = $2
				WHERE id = $3
			`
			_, err = tx.ExecContext(ctx, sqlStatement, translation.Language, translation.Value, translation.Id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
