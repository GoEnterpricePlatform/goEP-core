package mocksdata

import "github.com/GoEnterpricePlatform/goEP-core/addons/web-app/domain"

var Posts = []*domain.Post{
	{
		ID:      "1",
		Title:   "first post",
		Content: "first post with content",
	},
	{
		ID:      "2",
		Title:   "second post",
		Content: "second post with content",
	},
}