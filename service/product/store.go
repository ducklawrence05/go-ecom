package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ducklawrence/go-ecom/db"
	"github.com/ducklawrence/go-ecom/types"
)

type Store struct {
	db db.DBTX
}

func NewStore(db db.DBTX) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts(ctx context.Context) ([]types.Product, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM Products")
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		product, err := scanRowsInProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (s *Store) GetProductsByIDs(ctx context.Context, productIDs []int) ([]types.Product, error) {
	placeholders := make([]string, len(productIDs))
	args := make([]any, len(productIDs))

	for i, v := range productIDs {
		placeholders[i] = fmt.Sprintf("@p%d", i+1)
		args[i] = v
	}

	query := fmt.Sprintf(
		"SELECT * FROM Products WHERE ID IN (%s)",
		strings.Join(placeholders, ", "),
	)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0, len(productIDs))
	for rows.Next() {
		product, err := scanRowsInProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func scanRowsInProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) CreateProduct(ctx context.Context, product types.Product) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO Products (Name, Description, Image, Price, Quantity)
			VALUES (@Name, @Description, @Image, @Price, @Quantity)`,
		sql.Named("Name", product.Name),
		sql.Named("Description", product.Description),
		sql.Named("Image", product.Image),
		sql.Named("Price", product.Price),
		sql.Named("Quantity", product.Quantity),
	)

	return err
}

func (s *Store) UpdateProduct(ctx context.Context, product types.Product) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE Products SET Name = @Name, Price = @Price, Image = @Image, 
			Description = @Description, Quantity = @Quantity
			WHERE ID = @ID`,
		sql.Named("Name", product.Name),
		sql.Named("Price", product.Price),
		sql.Named("Image", product.Image),
		sql.Named("Description", product.Description),
		sql.Named("Quantity", product.Quantity),
		sql.Named("ID", product.ID),
	)

	return err
}
