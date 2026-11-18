package ai

import (
	"github.com/GoEnterpricePlatform/goEP-core/addons/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/ai/contract"
)

var _ contract.ToolProvider = &Provider{}

type Provider struct {
	PostSrv port.PostSrv
}

func NewPostAiProvider(postSrv port.PostSrv) *Provider {
	return &Provider{
		PostSrv: postSrv,
	}
}
