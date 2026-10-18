package openai

import (
	"fmt"

	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func save() {

	tools := []responses.ToolUnionParam{
		{
			OfFunction: &responses.FunctionToolParam{
				Name:        "create_role",
				Description: param.Opt[string]{Value: "create a new role in the system."},
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"title": map[string]any{
							"type":        "string",
							"description": "That's the title of the post.",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "This is the post description.",
						},
					},
					"required": []string{"title", "desc"},
				},
			},
		},
		{
			OfFunction: &responses.FunctionToolParam{
				Name:        "create_post",
				Description: param.Opt[string]{Value: "Create a new post in the system."},
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"title": map[string]any{
							"type":        "string",
							"description": "Title of the post.",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "Description of the post.",
						},
					},
					"required": []string{"title"},
				},
			},
		},
		{
			OfFunction: &responses.FunctionToolParam{
				Name:        "update_post",
				Description: param.Opt[string]{Value: "Update an existing post in the system."},
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id": map[string]any{
							"type":        "string",
							"description": "ID of the post to update.",
						},
						"title": map[string]any{
							"type":        "string",
							"description": "Updated title of the post.",
						},
						"desc": map[string]any{
							"type":        "string",
							"description": "Updated description of the post.",
						},
					},
					"required": []string{"id"},
				},
			},
		},
		{
			OfFunction: &responses.FunctionToolParam{
				Name:        "delete_post",
				Description: param.Opt[string]{Value: "Delete a post from the system."},
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id": map[string]any{
							"type":        "string",
							"description": "ID of the post to delete.",
						},
					},
					"required": []string{"id"},
				},
			},
		},
		{
			OfFunction: &responses.FunctionToolParam{
				Name:        "get_all_posts",
				Description: param.Opt[string]{Value: "Retrieve all posts from the system."},
				Parameters: map[string]any{
					"type":       "object",
					"properties": map[string]any{},
				},
			},
		},
	}
	fmt.Println(tools)
}
