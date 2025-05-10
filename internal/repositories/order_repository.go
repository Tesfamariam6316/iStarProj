package repositories

import (
	"context"
	"github.com/hulupay/istar-api/internal/models"
	"go.uber.org/zap"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	UpdateOrderStatus(ctx context.Context, orderID string, status models.OrderStatus, txHash *string, completedAt *time.Time, errorMessage *string) error
}

type orderRepository struct {
	/*db     *pgxpool.Pool*/
	logger *zap.Logger
}

func NewOrderRepository( /*db *pgxpool.Pool,*/ logger *zap.Logger) OrderRepository {
	return &orderRepository{ /*db: db,*/ logger: logger.Named("order_repository")}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	//query := `
	//	INSERT INTO orders (id, type, status, username, recipient_hash, quantity, months, amount, wallet_type, created_at, updated_at)
	//	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	//`
	//_, err := r.db.Exec(ctx, query,
	//	order.ID, order.Type, order.Status, order.Username, order.RecipientHash,
	//	order.Quantity, order.Months, order.Amount, order.WalletType,
	//	order.CreatedAt, order.UpdatedAt,
	//)
	//if err != nil {
	//	r.logger.Error("Failed to create order", zap.Error(err), zap.String("order_id", order.ID))
	//	return err
	//}
	return nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status models.OrderStatus, txHash *string, completedAt *time.Time, errorMessage *string) error {
	//query := `
	//	UPDATE orders
	//	SET status = $1, tx_hash = $2, completed_at = $3, error_message = $4, updated_at = $5
	//	WHERE id = $6
	//`
	//_, err := r.db.Exec(ctx, query, status, txHash, completedAt, errorMessage, time.Now(), orderID)
	//if err != nil {
	//	r.logger.Error("Failed to update order status", zap.Error(err), zap.String("order_id", orderID))
	//	return err
	//}
	return nil
}
