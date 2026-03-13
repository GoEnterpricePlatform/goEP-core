package ai

import (
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
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
