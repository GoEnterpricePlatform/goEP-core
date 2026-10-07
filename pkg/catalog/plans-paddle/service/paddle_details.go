package service

import (
	"context"

	paddlePlanD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	paddleSDK "github.com/PaddleHQ/paddle-go-sdk/v5"
)

// enrichPaddleDetails uses Paddle's product endpoint with its prices included,
// so each plan needs only one Paddle request regardless of its item count.
// Paddle data is optional for the catalog response: deleted or unavailable
// Paddle records leave their local IDs intact and the plan remains displayable.
func (s *Service) enrichPaddleDetails(ctx context.Context, plan *paddlePlanD.PaddlePlan) {
	if s.PaddleClient == nil || plan == nil || plan.PaddleProductID == "" {
		return
	}

	product, err := s.PaddleClient.GetProduct(ctx, &paddleSDK.GetProductRequest{
		ProductID:     plan.PaddleProductID,
		IncludePrices: true,
	})
	if err != nil || product == nil {
		return
	}

	plan.PaddleProduct = &paddlePlanD.PaddleProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		ImageURL:    product.ImageURL,
		Status:      string(product.Status),
	}

	prices := make(map[string]*paddleSDK.Price, len(product.Prices))
	for i := range product.Prices {
		prices[product.Prices[i].ID] = &product.Prices[i]
	}
	for _, item := range plan.Items {
		price, ok := prices[item.PaddlePriceID]
		if !ok {
			continue
		}

		item.PaddlePrice = &paddlePlanD.PaddlePrice{
			ID:           price.ID,
			Name:         price.Name,
			Description:  price.Description,
			Amount:       price.UnitPrice.Amount,
			CurrencyCode: string(price.UnitPrice.CurrencyCode),
			Status:       string(price.Status),
		}
		if price.BillingCycle != nil {
			item.PaddlePrice.BillingCycle = &paddlePlanD.PaddleBillingCycle{
				Interval:  string(price.BillingCycle.Interval),
				Frequency: price.BillingCycle.Frequency,
			}
		}
	}
}
