package service

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func validateLemonSqueezyCredentials(apiKey string) error {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.lemonsqueezy.com/v1/stores",
		nil,
	)
	if err != nil {
		return fmt.Errorf("creating lemonsqueezy request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.api+json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling lemonsqueezy api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read lemonsqueezy error response: %v", err)
		}
		return fmt.Errorf("%w", domain.ErrInvalidLemonSqueezyCredentials)
	}

	return nil
}
