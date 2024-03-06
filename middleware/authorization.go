package middleware

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
)

func Authorized(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelId := r.PathValue("modelId")
		email := r.Context().Value(UserIdentifierKey)

		sqlStatement := `
			SELECT COUNT(*)
			FROM model_user_relations
			WHERE modelId = $1 AND userId = (SELECT id FROM users WHERE email = $2)`

		var count int
		err := db.QueryRowContext(r.Context(), sqlStatement, modelId, email).Scan(&count)
		if err != nil {
			slog.WarnContext(r.Context(), fmt.Sprintf("error error checking if user is authorized for model ID %v: %v", modelId, err.Error()))
			http.Error(w, err.Error(), 500)
			return
		}

		if count == 0 {
			slog.InfoContext(r.Context(), fmt.Sprintf("user is not authorized for model ID %v", modelId))
			w.WriteHeader(http.StatusForbidden)
		}

		next(w, r)
	}
}
