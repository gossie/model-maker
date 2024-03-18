package persistence

import (
	"context"
	"database/sql"

	"github.com/gossie/modelling-service/domain"
)

type psqlConstraintRepository struct {
	db *sql.DB
}

func NewPsqlConstraintRepository(db *sql.DB) psqlConstraintRepository {
	return psqlConstraintRepository{db: db}
}

func (repo *psqlConstraintRepository) SaveConstraint(ctx context.Context, modelId string, ccr domain.ConstraintCreationRequest) (int, error) {
	var parameterId int
	sql := `
		INSERT INTO constraints (constraintType, fromId, fromValueId, targetId, targetValueId, modelId)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	err := repo.db.QueryRowContext(ctx, sql, ccr.Type, ccr.FromId, ccr.FromValueId, ccr.TargetId, ccr.TargetValueId, modelId).Scan(&parameterId)
	return parameterId, err
}

func (repo *psqlConstraintRepository) DeleteConstraint(ctx context.Context, modelId, constraintId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM constraints WHERE id = $1 AND modelId = $2", constraintId, modelId)
	return err
}
