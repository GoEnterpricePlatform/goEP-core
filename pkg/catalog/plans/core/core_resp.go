package core

import "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"

type PlanResp struct {
	Count    int64          `json:"count"`
	Pages    int64          `json:"pages"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Products []*domain.Plan `json:"plans"`
}
