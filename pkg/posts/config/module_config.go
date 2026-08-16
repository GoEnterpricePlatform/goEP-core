package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	postHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/handler"
	postP "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	postRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/repository/mongo"
	postService "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/service"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppStack   *config.AppStack
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB *mongo.Database
}

type Module struct {
	PostService   postP.PostSrv
}

func NewPostModule(cfg ModuleConfig) (*Module, error) {

	// module name
	mdlName := "posts"

	// collections
	postCollName := fmt.Sprintf("%s_posts", mdlName)
	postColl := cfg.DB.Collection(postCollName)

	// Repositories
	postRepo := postRepository.NewPostRepo(cfg.AppClients.MongoConn.DB, postColl)

	// Indexes
	err := postRepo.CreateIndexes()
	if err != nil {
		return nil, err
	}

	// services
	postSrv := postService.NewPostSrv(postRepo)

	// register handlers
	postHandler.NewPostApiHandler(cfg.APIv1, postSrv)

	return &Module{
		PostService: postSrv,
	}, nil
}
