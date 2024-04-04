package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	configurationmodel "github.com/gossie/configuration-model"
	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
)

type psqlModelRepository struct {
	db *sql.DB
}

func NewPsqlModelRepository(db *sql.DB) psqlModelRepository {
	return psqlModelRepository{db: db}
}

func (mr *psqlModelRepository) FindById(ctx context.Context, modelId int) (domain.Model, error) {
	sqlStatement := `
		SELECT m.id, m.name, t.translation, c.id, c.constraintType, c.fromId, c.fromValueId, c.targetId, c.targetValueId FROM models m
		LEFT JOIN model_translations t
		ON m.id = t.modelId
		LEFT JOIN constraints c
		ON c.modelId = m.id
		WHERE m.id = $1 AND t.language = $2
	`
	rows, err := mr.db.QueryContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey))
	if err != nil {
		slog.WarnContext(ctx, "got an error"+err.Error())
		return domain.Model{}, err
	}
	defer rows.Close()

	var id int
	var name string
	var translation sql.NullString
	constraints := make([]domain.Constraint, 0)
	for rows.Next() {
		var constraintId sql.NullInt32
		var contraintType sql.NullInt32
		var fromId sql.NullInt32
		var fromValueId sql.NullInt32
		var targetId sql.NullInt32
		var targetValueId sql.NullInt32

		rows.Scan(&id, &name, &translation, &constraintId, &contraintType, &fromId, &fromValueId, &targetId, &targetValueId)

		if constraintId.Valid {
			constraints = append(constraints, domain.Constraint{
				Id:            int(constraintId.Int32),
				Type:          configurationmodel.ConstraintType(contraintType.Int32),
				FromId:        int(fromId.Int32),
				FromValueId:   int(fromValueId.Int32),
				TargetId:      int(targetId.Int32),
				TargetValueId: int(targetValueId.Int32),
			})
		}
	}

	return domain.Model{Id: id, Name: name, Translation: translation.String, Constraints: constraints}, rows.Err()
}

func (mr *psqlModelRepository) FindAllByUser(ctx context.Context, userEmail string) ([]domain.Model, error) {
	language := ctx.Value(middleware.LanguageKey)

	slog.InfoContext(ctx, fmt.Sprintf("retrieving models for user %v and language %v", userEmail, language))

	sqlStatement := `
		SELECT m.id, m.name, t.translation
		FROM models m
		LEFT JOIN model_user_relations mur
		ON m.id = mur.modelid
		LEFT JOIN model_translations t
		ON m.id = t.modelId
		WHERE mur.userid = (SELECT id FROM users WHERE email = $1) AND (t.language = $2 OR t.language IS NULL)
	`
	rows, err := mr.db.QueryContext(ctx, sqlStatement, userEmail, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]domain.Model, 0)
	for rows.Next() {
		var id int
		var name string
		var translation sql.NullString
		rows.Scan(&id, &name, &translation)
		models = append(models, domain.Model{Id: id, Name: name, Translation: translation.String, Constraints: make([]domain.Constraint, 0)})
	}

	return models, rows.Err()
}

func (mr *psqlModelRepository) SaveModel(ctx context.Context, userEmail string, cmr domain.ModelCreationRequest) (int, error) {
	tx, _ := mr.db.BeginTx(ctx, nil)

	var userId int
	row := tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = $1", userEmail)
	err := row.Scan(&userId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}

	var modelId int
	err = mr.db.QueryRowContext(ctx, "INSERT INTO models (name) VALUES ($1) RETURNING id", cmr.Name).Scan(&modelId)
	if err != nil {
		_ = tx.Rollback()
		return -1, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("created model with ID %v", modelId))

	_, err = mr.db.ExecContext(ctx, "INSERT INTO model_user_relations (modelId, userId) VALUES ($1, $2)", modelId, userId)
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
