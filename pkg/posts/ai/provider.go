package ai

import "github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"

var _ contract.ToolProvider = &Provider{}

type Provider struct{}

func NewPostAiProvider() *Provider {
	return &Provider{}
}

func (p *Provider) GetAITools() []contract.ToolDefinition {
	return []contract.ToolDefinition{
		{
			Name:        "create_post",
			Description: "Create a new post in the system.",
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Title of the post",
					},
					"desc": map[string]interface{}{
						"type":        "string",
						"description": "Description of the post",
					},
					"author": map[string]interface{}{
						"type":        "string",
						"description": "author for this post",
					},
				},
				"required": []string{"title", "author"},
			},
		},
		{
			Name:        "update_post",
			Description: "Update an post in the system.",
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "ID of the post to update.",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Updated title of the post.",
					},
					"desc": map[string]interface{}{
						"type":        "string",
						"description": "Updated description of the post.",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "delete_post",
			Description: "Delete a post from the system.",
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "ID of the post to delete.",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			Name:        "get_all_posts",
			Description: "Retrieve all posts from the system.",
			Schemma: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}
