package paddle

import (
	paddle "github.com/PaddleHQ/paddle-go-sdk/v5"
)

type PaddleClient struct {
	Client *paddle.SDK
}

func NewPaddleClient(apiKey string, environment string) (*PaddleClient, error) {
	baseURL := paddle.SandboxBaseURL
	if environment == "production" {
		baseURL = paddle.ProductionBaseURL
	}
	client, err := paddle.New(apiKey, paddle.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}
	return &PaddleClient{
		Client: client,
	}, nil
}
