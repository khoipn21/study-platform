package lemonsqueezy

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/study-platform/payment-service/internal/model"
)

const (
	BaseURL = "https://api.lemonsqueezy.com/v1"
	TestBaseURL = "https://api.lemonsqueezy.com/v1" // Lemon Squeezy uses the same URL for test mode
)

type Client struct {
	apiKey     string
	storeID    string
	baseURL    string
	httpClient *http.Client
	webhookSecret string
}

type Config struct {
	APIKey        string
	StoreID       string
	Environment   string // "test" or "production"
	WebhookSecret string
}

// API Response structures following JSON:API specification
type APIResponse struct {
	Data   interface{}            `json:"data"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
	Errors []APIError             `json:"errors,omitempty"`
}

type APIError struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source struct {
		Pointer   string `json:"pointer,omitempty"`
		Parameter string `json:"parameter,omitempty"`
	} `json:"source,omitempty"`
}

type CheckoutData struct {
	Type       string            `json:"type"`
	ID         string            `json:"id,omitempty"`
	Attributes CheckoutAttributes `json:"attributes"`
	Relationships CheckoutRelationships `json:"relationships"`
}

type CheckoutAttributes struct {
	StoreID         int                    `json:"store_id,omitempty"`
	VariantID       int                    `json:"variant_id,omitempty"`
	Custom          map[string]interface{} `json:"custom,omitempty"`
	CheckoutOptions CheckoutOptions        `json:"checkout_options,omitempty"`
	CheckoutData    CheckoutDataOptions    `json:"checkout_data,omitempty"`
	Preview         bool                   `json:"preview,omitempty"`
	URL             string                 `json:"url,omitempty"`
}

type CheckoutOptions struct {
	Embed               bool   `json:"embed,omitempty"`
	Media               bool   `json:"media,omitempty"`
	Logo                bool   `json:"logo,omitempty"`
	Desc                bool   `json:"desc,omitempty"`
	Discount            bool   `json:"discount,omitempty"`
	ButtonColor         string `json:"button_color,omitempty"`
}

type CheckoutDataOptions struct {
	Email           string                 `json:"email,omitempty"`
	Name            string                 `json:"name,omitempty"`
	BillingAddress  map[string]string      `json:"billing_address,omitempty"`
	TaxNumber       string                 `json:"tax_number,omitempty"`
	DiscountCode    string                 `json:"discount_code,omitempty"`
	Custom          map[string]interface{} `json:"custom,omitempty"`
}

type CheckoutRelationships struct {
	Store struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"store"`
	Variant struct {
		Data struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"data"`
	} `json:"variant"`
}

func NewClient(config Config) *Client {
	baseURL := BaseURL
	if config.Environment == "test" {
		// Lemon Squeezy uses the same URL but different store/product IDs for testing
		baseURL = TestBaseURL
	}

	return &Client{
		apiKey:        config.APIKey,
		storeID:       config.StoreID,
		baseURL:       baseURL,
		webhookSecret: config.WebhookSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) CreateCheckout(ctx context.Context, req *model.LemonSqueezyCheckoutRequest) (*model.LemonSqueezyCheckoutResponse, error) {
	storeIDInt, err := strconv.Atoi(c.storeID)
	if err != nil {
		return nil, fmt.Errorf("invalid store ID: %w", err)
	}

	variantIDInt, err := strconv.Atoi(req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("invalid variant ID: %w", err)
	}

	checkoutData := CheckoutData{
		Type: "checkouts",
		Attributes: CheckoutAttributes{
			StoreID:   storeIDInt,
			VariantID: variantIDInt,
			Custom:    req.CustomData,
			CheckoutOptions: CheckoutOptions{
				Embed: true,
				Media: false,
				Logo:  true,
				Desc:  true,
			},
			CheckoutData: CheckoutDataOptions{
				Custom: req.CustomData,
			},
		},
		Relationships: CheckoutRelationships{
			Store: struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			}{
				Data: struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				}{
					Type: "stores",
					ID:   c.storeID,
				},
			},
			Variant: struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			}{
				Data: struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				}{
					Type: "variants",
					ID:   req.VariantID,
				},
			},
		},
	}

	payload := map[string]interface{}{
		"data": checkoutData,
	}

	respData, err := c.makeRequest(ctx, "POST", "/checkouts", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout: %w", err)
	}

	// Parse response
	var response struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				URL string `json:"url"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse checkout response: %w", err)
	}

	return &model.LemonSqueezyCheckoutResponse{
		CheckoutURL: response.Data.Attributes.URL,
		CheckoutID:  response.Data.ID,
	}, nil
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (*model.LemonSqueezyOrderData, error) {
	respData, err := c.makeRequest(ctx, "GET", fmt.Sprintf("/orders/%s", orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var response struct {
		Data model.LemonSqueezyOrderData `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse order response: %w", err)
	}

	return &response.Data, nil
}

func (c *Client) ListProducts(ctx context.Context) ([]interface{}, error) {
	respData, err := c.makeRequest(ctx, "GET", "/products", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	var response struct {
		Data []interface{} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse products response: %w", err)
	}

	return response.Data, nil
}

func (c *Client) ListVariants(ctx context.Context, productID string) ([]interface{}, error) {
	url := "/variants"
	if productID != "" {
		url += fmt.Sprintf("?filter[product_id]=%s", productID)
	}

	respData, err := c.makeRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	var response struct {
		Data []interface{} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse variants response: %w", err)
	}

	return response.Data, nil
}

// Variant represents a Lemon Squeezy variant
type Variant struct {
	ID         string  `json:"id"`
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	FormattedPrice string `json:"formatted_price"`
}

// GetVariants returns all available variants
func (c *Client) GetVariants() ([]Variant, error) {
	ctx := context.Background()
	respData, err := c.makeRequest(ctx, "GET", "/variants", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get variants: %w", err)
	}

	var response struct {
		Data []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				ProductID      int     `json:"product_id"`
				Name           string  `json:"name"`
				Price          int     `json:"price"` // Price in cents
				FormattedPrice string  `json:"formatted_price"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse variants response: %w", err)
	}

	var variants []Variant
	for _, v := range response.Data {
		variants = append(variants, Variant{
			ID:             v.ID,
			ProductID:      strconv.Itoa(v.Attributes.ProductID),
			Name:           v.Attributes.Name,
			Price:          float64(v.Attributes.Price) / 100.0, // Convert from cents
			FormattedPrice: v.Attributes.FormattedPrice,
		})
	}

	return variants, nil
}

func (c *Client) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) ([]byte, error) {
	url := c.baseURL + endpoint

	var reqBody []byte
	var err error

	if payload != nil {
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers according to Lemon Squeezy API requirements
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			respBody = append(respBody, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr struct {
			Errors []APIError `json:"errors"`
		}

		if json.Unmarshal(respBody, &apiErr) == nil && len(apiErr.Errors) > 0 {
			return nil, fmt.Errorf("API error: %s - %s", apiErr.Errors[0].Title, apiErr.Errors[0].Detail)
		}

		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// VerifyWebhookSignature verifies that a webhook payload was sent by Lemon Squeezy
func (c *Client) VerifyWebhookSignature(payload []byte, signature string) bool {
	if c.webhookSecret == "" {
		return false
	}

	hash := hmac.New(sha256.New, []byte(c.webhookSecret))
	hash.Write(payload)
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// ParseWebhookPayload parses a webhook payload from Lemon Squeezy
func (c *Client) ParseWebhookPayload(payload []byte) (*model.LemonSqueezyWebhookPayload, error) {
	var webhookPayload model.LemonSqueezyWebhookPayload

	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	return &webhookPayload, nil
}