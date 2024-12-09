package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents an HTTP client for the barn_manager service
type bClient struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new client for the barn_manager service
func NewClient(baseURL string) *bClient {
	return &bClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

var JsonMarshal = json.Marshal

// GetProducts retrieves products from the barn_manager service
func (c *bClient) GetProducts(ctx context.Context, recipes []string) (string, error) {

	data, err := JsonMarshal(recipes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal product: %w", err)
	}

	resp, err := c.client.Post(c.baseURL+"/api/v1/products", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to get products: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var products string
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return products, nil
}
