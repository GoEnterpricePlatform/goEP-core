package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func validateStripeCredentials(secretKey string) error {
	req, err := http.NewRequest(http.MethodGet, "https://api.stripe.com/v1/balance", nil)
	if err != nil {
		return fmt.Errorf("creating stripe request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling stripe api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read stripe error response: %v",err)
		}
		return fmt.Errorf("%w: %s", domain.ErrInvalidStripeCredentials, string(body))
	}

	return nil
}
