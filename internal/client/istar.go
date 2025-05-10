package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hulupay/istar-api/config"
	"github.com/hulupay/istar-api/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type IStarClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewIStarClient(cfg config.IStarConfig, logger *zap.Logger) *IStarClient {
	return &IStarClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
		},
		logger: logger.Named("istar_client"),
	}
}

func (c *IStarClient) DoRequest(ctx context.Context, method, path string, payload []byte) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		c.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to send request", zap.Error(err))
		return nil, fmt.Errorf("sending request failed: %w", err)
	}
	return resp, nil
}

func (c *IStarClient) CreateStarOrderAsync(ctx context.Context, req models.CreateStarOrderRequest) (*models.StarOrderResponse, error) {
	path := "/orders/star"
	payload, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return nil, models.InternalServerError("Failed to marshal request")
	}

	resp, err := c.DoRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, models.ValidationError("Invalid request parameters")
		case http.StatusUnauthorized:
			return nil, models.UnauthorizedError("Invalid API key")
		case http.StatusNotFound:
			return nil, models.NotFoundError("Resource not found")
		default:
			return nil, models.InternalServerError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		}
	}

	var response models.StarOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode response", zap.Error(err))
		return nil, models.InternalServerError("Failed to decode response")
	}

	c.logger.Info("Star order created (async)", zap.String("order_id", response.OrderID))
	return &response, nil
}

func (c *IStarClient) CreateStarOrderSync(ctx context.Context, req models.CreateStarOrderRequest) (*models.StarOrderResponse, error) {
	path := "/orders/star/sync"
	payload, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return nil, models.InternalServerError("Failed to marshal request")
	}

	resp, err := c.DoRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, models.ValidationError("Invalid request parameters")
		case http.StatusUnauthorized:
			return nil, models.UnauthorizedError("Invalid API key")
		case http.StatusNotFound:
			return nil, models.NotFoundError("Resource not found")
		default:
			return nil, models.InternalServerError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		}
	}

	var response models.StarOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode response", zap.Error(err))
		return nil, models.InternalServerError("Failed to decode response")
	}

	c.logger.Info("Star order created (sync)", zap.String("order_id", response.OrderID))
	return &response, nil
}

func (c *IStarClient) CreatePremiumOrderAsync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.PremiumOrderResponse, error) {
	path := "/orders/premium"
	payload, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return nil, models.InternalServerError("Failed to marshal request")
	}

	resp, err := c.DoRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, models.ValidationError("Invalid request parameters")
		case http.StatusUnauthorized:
			return nil, models.UnauthorizedError("Invalid API key")
		case http.StatusNotFound:
			return nil, models.NotFoundError("Resource not found")
		default:
			return nil, models.InternalServerError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		}
	}

	var response models.PremiumOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode response", zap.Error(err))
		return nil, models.InternalServerError("Failed to decode response")
	}

	c.logger.Info("Premium order created (async)", zap.String("order_id", response.OrderID))
	return &response, nil
}

func (c *IStarClient) CreatePremiumOrderSync(ctx context.Context, req models.CreatePremiumOrderRequest) (*models.PremiumOrderResponse, error) {
	path := "/orders/premium/sync"
	payload, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal request", zap.Error(err))
		return nil, models.InternalServerError("Failed to marshal request")
	}

	resp, err := c.DoRequest(ctx, "POST", path, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode), zap.String("body", string(body)))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return nil, models.ValidationError("Invalid request parameters")
		case http.StatusUnauthorized:
			return nil, models.UnauthorizedError("Invalid API key")
		case http.StatusNotFound:
			return nil, models.NotFoundError("Resource not found")
		default:
			return nil, models.InternalServerError(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		}
	}

	var response models.PremiumOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode response", zap.Error(err))
		return nil, models.InternalServerError("Failed to decode response")
	}

	c.logger.Info("Premium order created (sync)", zap.String("order_id", response.OrderID))
	return &response, nil
}
