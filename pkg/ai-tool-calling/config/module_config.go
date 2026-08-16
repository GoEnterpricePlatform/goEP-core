package config

import (
	"log"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"

	tcP "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/port"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/adapter/disabled"
	openaiAdapter "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/adapter/open-ai"
	aiTCItz "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/initializer"
	toolCallingService "github.com/GoEnterpricePlatform/goEP-core/pkg/ai-tool-calling/service"
	postAI "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/ai"
	postModule "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppStack   *config.AppStack
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	// Module dependencies
	PostModule *postModule.Module

	DB *mongo.Database
}

type Module struct {
	TCService tcP.ToolCallingSrv
}

func NewAiModule(cfg ModuleConfig) (*Module, error) {

	var toolCallingAdapter tcP.ToolCallingAdt

	// The openai service is optional, and it is the responsibility of the developer or admin to configure it
	// this way we don't touch the code, it would only be from the UI to see how to handle it
	var openAiAdt *openaiAdapter.Adapter
	var disabledAdt *disabled.DisabledAdapter
	if cfg.AppClients.OpenaiCli != nil {
		openAiAdt = openaiAdapter.NewToolCallingAdt(
			cfg.AppClients.OpenaiCli.Client,
		)
		toolCallingAdapter = openAiAdt
	} else {
		disabledAdt = disabled.NewDisabledAdapter()
		toolCallingAdapter = disabledAdt
	}

	// AI - initializer
	aiTCInitializer := aiTCItz.NewAIItz()

	// AI - system-prompt
	systemPrompt, err := aiTCInitializer.GetSystemPrompt()
	if err != nil {
		log.Fatal(err)
	}

	toolCallingSrv := toolCallingService.NewToolCallingSrv(toolCallingAdapter, aiTCInitializer, systemPrompt)

	// AI - Providers
	postAIprovider := postAI.NewPostAiProvider(cfg.PostModule.PostService)

	// AI - register
	aiTCInitializer.RegisterTool(postAIprovider)

	// Set Tools Beofore register Register ALl modules
	if openAiAdt != nil {
		openAiAdt.SetTools(aiTCInitializer)
	} else {
		toolCallingAdapter = disabled.NewDisabledAdapter()
	}

	return &Module{
		TCService: toolCallingSrv,
	}, nil
}
