package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gossie/modelling-service/domain"
	"github.com/gossie/modelling-service/middleware"
)

type psqlModelRepository struct {
	db *sql.DB
}

func NewPsqlModelRepository(db *sql.DB) psqlModelRepository {
	return psqlModelRepository{db: db}
}

func (mr *psqlModelRepository) FindById(ctx context.Context, modelId string) (domain.Model, error) {
	sqlStatement := `
		SELECT m.id, m.name, t.translation FROM models m
		JOIN model_translations t
		ON m.id = t.modelId
		WHERE m.id = $1 AND t.language = $2
	`
	var id int
	var name string
	var translation sql.NullString
	row := mr.db.QueryRowContext(ctx, sqlStatement, modelId, ctx.Value(middleware.LanguageKey))
	err := row.Scan(&id, &name, &translation)

	return domain.Model{Id: id, Name: name, Translation: translation.String}, err
}

func (mr *psqlModelRepository) FindAllByUser(ctx context.Context, userEmail string) ([]domain.Model, error) {
	sqlStatement := `
		SELECT m.id, m.name, t.translation
		FROM models m
		JOIN model_user_relations mur
		ON m.id = mur.modelid
		JOIN model_translations t ON m.id = t.modelId
		WHERE mur.userid = (SELECT id FROM users WHERE email = $1) AND t.language = $2;
	`
	rows, err := mr.db.QueryContext(ctx, sqlStatement, userEmail, ctx.Value(middleware.LanguageKey))
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
		models = append(models, domain.Model{Id: id, Name: name, Translation: translation.String})
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
