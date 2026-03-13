package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/ai/contract"
)

func (p *Provider) GetAITools() []contract.ToolDefinition {
	return []contract.ToolDefinition{
		{
			Name:        "create_post",
			Description: "Create a new post in the system.",
			Category:    contract.ToolCategoryWrite,
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"prepare", "execute"},
						"description": "prepare = collect or modify fields, execute = perform the operation",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Title of the post",
					},
					"desc": map[string]interface{}{
						"type":        "string",
						"description": "Description of the post",
					},
				},
				"required": []string{"title"},
			},
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				title := args["title"].(string)

				// "desc" can have three states:
				// 1. Not present in args map  -> descPtr remains nil (field not provided)
				// 2. Present but empty ""     -> descPtr points to "" (empty value)
				// 3. Present with a value     -> descPtr points to the provided string
				var descPtr *string
				if v, ok := args["desc"].(string); ok {
					descPtr = &v
				}

				req := &core.CreatePostReq{
					Title: title,
					Desc:  descPtr,
				}

				if err := req.Validate(); err != nil {
					return nil, err
				}

				post := &domain.Post{
					Title: req.Title,
					Desc:  req.Desc,
				}

				err := p.PostSrv.Create(ctx, post)
				if err != nil {
					return nil, err
				}

				resp := map[string]any{
					"columns": []string{"title", "desc"},
					"rows": []map[string]any{
						{
							"title": post.Title,
							"desc":  post.Desc,
						},
					},
				}
				return resp, nil
			},
		},
		{
			Name:        "update_post",
			Description: "Update an post in the system.",
			Category:    contract.ToolCategoryWrite,
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"action": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"prepare", "execute"},
						"description": "prepare = collect or modify fields, execute = perform the operation",
					},
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
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				id := args["id"].(string)

				var titlePtr *string
				if v, ok := args["title"].(string); ok {
					titlePtr = &v
				}

				var descPtr *string
				if v, ok := args["desc"].(string); ok {
					descPtr = &v
				}

				req := &core.PatchPostReq{
					Title: titlePtr,
					Desc:  descPtr,
				}

				if err := req.Validate(); err != nil {
					return nil, err
				}

				post := &domain.Post{
					Desc: req.Desc,
				}

				if req.Title != nil {
					post.Title = *req.Title
				}

				updated, err := p.PostSrv.Patch(context.Background(), id, post)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"columns": []string{"id", "title", "desc"},
					"rows": []map[string]any{
						{
							"id":    updated.ID,
							"title": updated.Title,
							"desc":  updated.Desc,
						},
					},
				}, nil
			},
		},
		{
			Name:        "delete_post",
			Description: "Delete a post from the system.",
			Category:    contract.ToolCategoryWrite,
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
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				id := args["id"].(string)

				if id == "" {
					return nil, errors.New("missing post id from handler")
				}

				err := p.PostSrv.Delete(ctx, id)
				if err != nil {
					return nil, err
				}
				return map[string]any{
					"columns": []string{"result", "id"},
					"rows": []map[string]any{
						{
							"result": "post deleted",
							"id":     id,
						},
					},
				}, nil
			},
		},
		{
			Name:        "get_all_posts",
			Description: "Retrieve all posts from the system.",
			Category:    contract.ToolCategoryRead,
			Schemma: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				fmt.Println("############################################7")

				posts, err := p.PostSrv.GetAll(ctx)
				if err != nil {
					return nil, err
				}
				rows := []map[string]any{}
				fmt.Println("############################################8")

				for _, p := range posts {
					rows = append(rows, map[string]any{
						"id":    p.ID,
						"title": p.Title,
						"desc":  p.Desc,
					})
				}
				fmt.Println("############################################9")

				return map[string]any{
					"columns": []string{"id", "title", "desc"},
					"rows":    rows,
				}, nil
			},
		},
	}
}
