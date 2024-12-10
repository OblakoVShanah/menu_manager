package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	common "menu_manager/internal/models"
	"net/http"
	"strings"
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

	recipeString := "[" + strings.Join(recipes, ", ") + "]"
	log.Println(recipeString)
	data := []byte(recipeString)

	// data, err := JsonMarshal(recipeString)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal product: %w", err)
	// }

	resp, err := c.client.Post(c.baseURL+"/api/v1/products/check-availability", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to get products: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var productResp struct {
		Products []common.Product `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Преобразуем результат в строку
	resultBytes, err := json.Marshal(productResp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal product response: %w", err)
	}
	return string(resultBytes), nil
}
