package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type MyError struct {
	Msg        string
	StatusCode int
}

func (e MyError) Error() string {
	return e.Msg
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WithTransaction(
	ctx context.Context,
	db *sql.DB,
	fn func(tx *sql.Tx) *MyError,
) *MyError {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return &MyError{
			Msg:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil
}
