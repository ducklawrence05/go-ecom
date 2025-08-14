package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ducklawrence/go-ecom/db"
	"github.com/ducklawrence/go-ecom/types"
)

type Store struct {
	db db.DBTX
}

func NewStore(db db.DBTX) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT * FROM Users WHERE Email = @Email",
		sql.Named("Email", email),
	)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) GetUserByID(ctx context.Context, id int) (*types.User, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT * FROM Users WHERE ID = @ID",
		sql.Named("ID", id),
	)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(ctx context.Context, user types.User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO Users (FirstName, LastName, Email, Password) 
			VALUES (@FirstName, @LastName, @Email, @Password)`,
		sql.Named("FirstName", user.FirstName),
		sql.Named("LastName", user.LastName),
		sql.Named("Email", user.Email),
		sql.Named("Password", user.Password),
	)

	return err
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
