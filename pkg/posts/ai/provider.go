package ai

import (
	"context"
	"errors"

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
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to show to the user",
						"message":     "",
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
			Name:        "get_post_by_id",
			Description: "Retrieve a single post by its ID.",
			Category:    contract.ToolCategoryRead,
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "ID of the post to retrieve.",
					},
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to show to the user",
						"message":     "",
					},
				},
				"required": []string{"id"},
			},
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				id, ok := args["id"].(string)
				if !ok || id == "" {
					return nil, errors.New("missing post id")
				}

				post, err := p.PostSrv.Get(ctx, id)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"columns": []string{"id", "title", "desc"},
					"rows": []map[string]any{
						{
							"id":    post.ID,
							"title": post.Title,
							"desc":  deref(post.Desc),
						},
					},
				}, nil
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
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to show to the user",
						"message":     "",
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
							"desc":  deref(updated.Desc),
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
				id, ok := args["id"].(string)
				if !ok || id == "" {
					return nil, errors.New("missing post id")
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
				posts, err := p.PostSrv.GetAll(ctx)
				if err != nil {
					return nil, err
				}
				rows := []map[string]any{}

				for _, p := range posts {
					rows = append(rows, map[string]any{
						"id":    p.ID,
						"title": p.Title,
						"desc":  p.Desc,
					})
				}
				return map[string]any{
					"columns": []string{"id", "title", "desc"},
					"rows":    rows,
				}, nil
			},
		},
		{
			Name:        "search_posts",
			Description: "Search posts by title or description.",
			Category:    contract.ToolCategoryWrite,
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Keyword to search in post title or description",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results",
					},
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to show to the user",
						"message":     "",
					},
				},
				"required": []string{"query"},
			},
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				query := args["query"].(string)

				limit := 5
				if v, ok := args["limit"].(int); ok {
					limit = v
				}

				if query == "" {
					return nil, errors.New("missing query parameter")
				}

				// fmt.Println(limit)
				//fmt.Println(query)

				posts, err := p.PostSrv.Search(context.Background(), query, limit)
				if err != nil {
					return nil, err
				}

				rows := []map[string]any{}

				for _, post := range posts {
					rows = append(rows, map[string]any{
						"id":    post.ID,
						"title": post.Title,
						"desc":  deref(post.Desc),
					})
				}

				return map[string]any{
					"columns": []string{"id", "title", "desc"},
					"rows":    rows,
				}, nil
			},
		},
		{
			Name:        "ask_user",
			Description: "Ask the user for missing information.",
			Schemma: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to show to the user",
						"example":     "Para continuar necesito identificar el post.",
						"message":     "",
					},
					"fields": map[string]interface{}{
						"type":        "array",
						"description": "Fields the user can provide",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"name": map[string]interface{}{
									"type":    "string",
									"example": "12345",
								},
								"description": map[string]interface{}{
									"type":    "string",
									"example": "Identificador único del post",
								},
								"example": map[string]interface{}{
									"type":    "string",
									"example": "12345",
								},
							},
							"required": []string{"name", "description"},
						},
						"example": []map[string]interface{}{
							{
								"name":        "id",
								"description": "Identificador único del post",
								"example":     "12345",
							},
							{
								"name":        "title",
								"description": "Título del post",
								"example":     "Create http server",
							},
							{
								"name":        "description",
								"description": "Palabra clave en la descripción",
								"example":     "servidor HTTP",
							},
						},
					},
				},
				"required": []string{"message", "fields"},
			},
			Handler: func(ctx context.Context, args map[string]any) (map[string]any, error) {
				msg, _ := args["message"].(string)
				fieldsRaw, _ := args["fields"].([]any)

				rows := []map[string]any{}

				for _, f := range fieldsRaw {
					field := f.(map[string]any)

					rows = append(rows, map[string]any{
						"field":       field["name"],
						"description": field["description"],
						"example":     field["example"],
					})
				}

				return map[string]any{
					"message": msg,
					"columns": []string{"field", "description", "example"},
					"rows":    rows,
				}, nil
			},
		},
	}
}

// deref converts a string pointer (*string) into a regular string.
// If the pointer is nil, it returns an empty string ("") to avoid panics.
// Useful when working with optional fields (pointers) and you want to
// expose or show the actual value instead of the memory address.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
