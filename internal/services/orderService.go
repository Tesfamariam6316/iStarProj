package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/hulupay/
	"github.com/hulupay/istar-api/internal/client"
	"github.com/hulupay/istar-api/internal/models"
	"github.com/hulupay/istar-api/internal/repositories"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

// OrderService defines the interface for order-related business logic
type OrderService interface {
	CreateStarOrderAsync(ctx context.Context, req models.CreateStarOrderRequest) (*models.Order, error)
	CreateStarOrderSync(ctx context.Context, req models.CreateStarOrderRequest) (*models.Order, error)
	CreatePremiumOrderAsync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.Order, error)
	CreatePremiumOrderSync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.Order, error)
}

// orderService implements the OrderService interface
type orderService struct {
	repo        repositories.OrderRepository
	istarClient *client.IStarClient
	logger      *zap.Logger
}

// NewOrderService initializes a new OrderService with dependencies
func NewOrderService(repo repositories.OrderRepository, istarClient *client.IStarClient, logger *zap.Logger) OrderService {
	return &orderService{
		repo:        repo,
		istarClient: istarClient,
		logger:      logger.Named("order_service"),
	}
}

// CreateStarOrderAsync creates an asynchronous star gift order
func (s *orderService) CreateStarOrderAsync(ctx context.Context, req models.CreateStarOrderRequest) (*models.Order, error) {
	resp, err := s.istarClient.CreateStarOrderAsync(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create star order via iStar API", zap.Error(err))
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		s.logger.Error("Failed to parse created_at", zap.Error(err))
		return nil, models.InternalServerError("Invalid created_at timestamp")
	}

	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		s.logger.Error("Invalid order_id from iStar", zap.Error(err))
		return nil, models.InternalServerError("Invalid order_id")
	}

	order := &models.Order{
		ID:            orderID,
		Type:          models.OrderTypeStar,
		Status:        models.StatusPending,
		Username:      req.Username,
		RecipientHash: req.RecipientHash,
		Quantity:      &resp.Quantity,
		Amount:        resp.Amount,
		WalletType:    req.WalletType,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to save order to database", zap.Error(err))
		return nil, models.InternalServerError("Failed to save order")
	}

	s.logger.Info("Star order created (async)", zap.String("order_id", order.ID.String()))
	return order, nil
}

// CreateStarOrderSync creates a synchronous star gift order
func (s *orderService) CreateStarOrderSync(ctx context.Context, req models.CreateStarOrderRequest) (*models.Order, error) {
	resp, err := s.istarClient.CreateStarOrderSync(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create star order via iStar API", zap.Error(err))
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		s.logger.Error("Failed to parse created_at", zap.Error(err))
		return nil, models.InternalServerError("Invalid created_at timestamp")
	}

	var completedAt *time.Time
	if resp.CompletedAt != nil {
		t, err := time.Parse(time.RFC3339, *resp.CompletedAt)
		if err != nil {
			s.logger.Error("Failed to parse completed_at", zap.Error(err))
			return nil, models.InternalServerError("Invalid completed_at timestamp")
		}
		completedAt = &t
	}

	status := models.OrderStatus(resp.Status)
	if status != models.StatusCompleted && status != models.StatusFailed {
		s.logger.Warn("Unexpected status from iStar", zap.String("status", resp.Status))
		status = models.StatusFailed
	}

	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		s.logger.Error("Invalid order_id from iStar", zap.Error(err))
		return nil, models.InternalServerError("Invalid order_id")
	}

	order := &models.Order{
		ID:            orderID,
		Type:          models.OrderTypeStar,
		Status:        status,
		Username:      req.Username,
		RecipientHash: req.RecipientHash,
		Quantity:      &resp.Quantity,
		Amount:        resp.Amount,
		WalletType:    req.WalletType,
		TxHash:        resp.TxHash,
		CreatedAt:     createdAt,
		UpdatedAt:     time.Now(),
		CompletedAt:   completedAt,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to save order to database", zap.Error(err))
		return nil, models.InternalServerError("Failed to save order")
	}

	s.logger.Info("Star order created (sync)", zap.String("order_id", order.ID.String()))
	return order, nil
}

// CreatePremiumOrderAsync creates an asynchronous premium gift order
func (s *orderService) CreatePremiumOrderAsync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.Order, error) {
	resp, err := s.istarClient.CreatePremiumOrderAsync(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create premium order via iStar API", zap.Error(err))
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		s.logger.Error("Failed to parse created_at", zap.Error(err))
		return nil, models.InternalServerError("Invalid created_at timestamp")
	}

	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		s.logger.Error("Invalid order_id from iStar", zap.Error(err))
		return nil, models.InternalServerError("Invalid order_id")
	}

	order := &models.Order{
		ID:            orderID,
		Type:          models.OrderTypePremium,
		Status:        models.StatusPending,
		Username:      req.Username,
		RecipientHash: req.RecipientHash,
		Months:        &resp.Months,
		Amount:        resp.Amount,
		WalletType:    req.WalletType,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to save order to database", zap.Error(err))
		return nil, models.InternalServerError("Failed to save order")
	}

	s.logger.Info("Premium order created (async)", zap.String("order_id", order.ID.String()))
	return order, nil
}

// CreatePremiumOrderSync creates a synchronous premium gift order
func (s *orderService) CreatePremiumOrderSync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.Order, error) {
	resp, err := s.istarClient.CreatePremiumOrderSync(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create premium order via iStar API", zap.Error(err))
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		s.logger.Error("Failed to parse created_at", zap.Error(err))
		return nil, models.InternalServerError("Invalid created_at timestamp")
	}

	var completedAt *time.Time
	if resp.CompletedAt != nil {
		t, err := time.Parse(time.RFC3339, *resp.CompletedAt)
		if err != nil {
			s.logger.Error("Failed to parse completed_at", zap.Error(err))
			return nil, models.InternalServerError("Invalid completed_at timestamp")
		}
		completedAt = &t
	}

	status := models.OrderStatus(resp.Status)
	if status != models.StatusCompleted && status != models.StatusFailed {
		s.logger.Warn("Unexpected status from iStar", zap.String("status", resp.Status))
		status = models.StatusFailed
	}

	orderID, err := uuid.Parse(resp.OrderID)
	if err != nil {
		s.logger.Error("Invalid order_id from iStar", zap.Error(err))
		return nil, models.InternalServerError("Invalid order_id")
	}

	order := &models.Order{
		ID:            orderID,
		Type:          models.OrderTypePremium,
		Status:        status,
		Username:      req.Username,
		RecipientHash: req.RecipientHash,
		Months:        &resp.Months,
		Amount:        resp.Amount,
		WalletType:    req.WalletType,
		TxHash:        resp.TxHash,
		CreatedAt:     createdAt,
		UpdatedAt:     time.Now(),
		CompletedAt:   completedAt,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.logger.Error("Failed to save order to database", zap.Error(err))
		return nil, models.InternalServerError("Failed to save order")
	}

	s.logger.Info("Premium order created (sync)", zap.String("order_id", order.ID.String()))
	return order, nil
}
