package order

import (
	"context"
	"database/sql"

	"github.com/ducklawrence/go-ecom/db"
	"github.com/ducklawrence/go-ecom/types"
)

type Store struct {
	db db.DBTX
}

func NewStore(db db.DBTX) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(ctx context.Context, order types.Order) (int, error) {
	var insertedID int
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO Orders (UserID, Total, Status, Address)
			OUTPUT INSERTED.ID
			VALUES (@UserID, @Total, @Status, @Address)`,
		sql.Named("UserID", order.UserID),
		sql.Named("Total", order.Total),
		sql.Named("Status", order.Status),
		sql.Named("Address", order.Address),
	).Scan(&insertedID)

	if err != nil {
		return 0, err
	}

	return int(insertedID), nil
}

func (s *Store) CreateOrderItem(ctx context.Context, orderItem types.OrderItem) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO OrderItems (OrderID, ProductID, Quantity, Price)
			VALUES (@OrderID, @ProductID, @Quantity, @Price)`,
		sql.Named("OrderID", orderItem.OrderID),
		sql.Named("ProductID", orderItem.ProductID),
		sql.Named("Quantity", orderItem.Quantity),
		sql.Named("Price", orderItem.Price),
	)
	return err
}
