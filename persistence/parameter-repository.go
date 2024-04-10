package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	configurationmodel "github.com/gossie/configuration-model"
	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
)

type psqlParameterRepository struct {
	db *sql.DB
}

func NewPsqlParameterRepository(db *sql.DB) psqlParameterRepository {
	return psqlParameterRepository{db: db}
}

func (pr *psqlParameterRepository) FindAllByModelId(ctx context.Context, modelId int, searchValue string) ([]domain.Parameter, error) {
	if searchValue == "*" {
		searchValue = ""
	}

	var (
		rows *sql.Rows
		err  error
	)

	if searchValue == "" {
		slog.InfoContext(ctx, fmt.Sprintf("retrieving all parameters of model with ID %v", modelId))

		sqlStatement := `
			SELECT p.id, p.name, p.valueType, pt.translation, v.id, v.value, vt.translation
			FROM parameters p
			LEFT JOIN parameter_translations pt
			ON p.id = pt.parameterId
			LEFT JOIN values v
			ON v.parameterId = p.id
			LEFT JOIN value_translations vt
			ON vt.valueId = v.id
			WHERE p.modelId = $1
			AND (pt.language = $2 OR pt.language IS NULL)
			AND (vt.language = $2 OR vt.language IS NULL)
			ORDER BY p.id
		`
		rows, err = pr.db.QueryContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey))
	} else {
		slog.InfoContext(ctx, fmt.Sprintf("searching for parameters containing '%v' at model with ID %v", searchValue, modelId))

		sqlStatement := `
			SELECT p.id, p.name, p.valueType, pt.translation, v.id, v.value, vt.translation
			FROM parameters p
			LEFT JOIN parameter_translations pt
			ON p.id = pt.parameterId
			LEFT JOIN values v
			ON v.parameterId = p.id
			LEFT JOIN value_translations vt
			ON vt.valueId = v.id
			WHERE p.modelId = $1
			AND (pt.language = $2 OR pt.language IS NULL)
			AND (vt.language = $2 OR vt.language IS NULL)
			AND (p.name LIKE '%' || $3 || '%' OR pt.translation LIKE '%' || $3 || '%')
			ORDER BY p.id
		`
		rows, err = pr.db.QueryContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey), searchValue)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parameters := make([]domain.Parameter, 0)
	for rows.Next() {
		var id int
		var name string
		var valueType configurationmodel.ValueType
		var paramTranslation sql.NullString
		var valueId sql.NullInt32
		var paramValue sql.NullString
		var valueTranslation sql.NullString
		err = rows.Scan(&id, &name, &valueType, &paramTranslation, &valueId, &paramValue, &valueTranslation)
		if err != nil {
			return nil, err
		}

		lastIndex := len(parameters) - 1
		if lastIndex < 0 || parameters[lastIndex].Id != id {
			parameters = append(parameters, domain.Parameter{
				Id:          id,
				Name:        name,
				Translation: paramTranslation.String,
				ValueType:   valueType,
				Value: domain.ParameterValue{
					Values: make([]domain.Value, 0),
				},
			})
			lastIndex++
		}

		if valueId.Valid {
			parameters[lastIndex].Value.Values = append(parameters[lastIndex].Value.Values, domain.Value{Id: int(valueId.Int32), Value: paramValue.String, Translation: valueTranslation.String})
		}
	}

	return parameters, nil
}

func (pr *psqlParameterRepository) SaveParameter(ctx context.Context, modelId int, pmr domain.ParameterCreationRequest) (int, error) {
	var parameterId int
	err := pr.db.QueryRowContext(ctx, "INSERT INTO parameters (name, valueType, modelId) VALUES ($1, $2, $3) RETURNING id", pmr.Name, pmr.ValueType, modelId).Scan(&parameterId)
	return parameterId, err
}

func (pr *psqlParameterRepository) DeleteParameter(ctx context.Context, modelId int, parameterId int) error {
	_, err := pr.db.ExecContext(ctx, "DELETE FROM parameters WHERE id = $1 AND modelId = $2", parameterId, modelId)
	return err
}

func (pr *psqlParameterRepository) FindAllTranslations(ctx context.Context, parameterId string) ([]domain.Translation, error) {
	sqlStatement := `
		SELECT id, field, language, translation
		FROM parameter_translations
		WHERE parameterId = $1
	`
	rows, err := pr.db.QueryContext(ctx, sqlStatement, parameterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	translations := make([]domain.Translation, 0)

	for rows.Next() {
		var id int
		var field string
		var language string
		var value string
		rows.Scan(&id, &field, &language, &value)

		translations = append(translations, domain.Translation{Id: id, Field: field, Language: language, Value: value})
	}

	return translations, rows.Err()
}

func (pr *psqlParameterRepository) SaveTranslations(ctx context.Context, parameterId string, tmr domain.TranslationModificationRequest) error {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if len(tmr.NewTranslations) > 0 {
		args := make([]any, 0, len(tmr.NewTranslations)*4)
		valueStrings := make([]string, 0, len(tmr.NewTranslations))
		for _, translation := range tmr.NewTranslations {
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

	if len(tmr.UpdatedTranslations) > 0 {
		for _, translation := range tmr.UpdatedTranslations {
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

func (pr *psqlParameterRepository) SaveValues(ctx context.Context, parameterId string, vmr domain.ValueModificationRequest) error {
	tx, err := pr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if len(vmr.NewValues) > 0 {
		args := make([]any, 0, len(vmr.NewValues)*2)
		valueStrings := make([]string, 0, len(vmr.NewValues))
		for _, value := range vmr.NewValues {
			valueStrings = append(valueStrings, fmt.Sprintf("($%v, $%v)", len(args)+1, len(args)+2))
			args = append(args, value, parameterId)
		}

		sqlStatement := `
			INSERT INTO values (value, parameterId)
			VALUES ` + strings.Join(valueStrings, ", ")
		_, err = tx.ExecContext(ctx, sqlStatement, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(vmr.UpdatedValues) > 0 {
		for _, value := range vmr.UpdatedValues {
			sqlStatement := `
				UPDATE values
				SET value = $1
				WHERE id = $2
			`
			_, err = tx.ExecContext(ctx, sqlStatement, value.Value, value.Id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
