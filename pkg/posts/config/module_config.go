package config

import (
	"fmt"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/internal/config"
	postHandler "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/handler"
	postP "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/port"
	postRepository "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/repository/mongo"
	postService "github.com/GoEnterpricePlatform/goEP-core/pkg/posts/service"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/middlewares"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModuleConfig struct {
	AppStack   *config.AppStack
	AppEnvs    *config.AppEnvs
	AppClients *config.AppClients
	APIv1      *http.ServeMux

	DB *mongo.Database

	Deps ModuleDeps
}

// ModuleDeps groups the external dependencies (from other modules/services)
// that this module needs to function.
type ModuleDeps struct {
	AuthApiMdw *middlewares.AuthMiddleware
}

type Module struct {
	PostService postP.PostSrv
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
	postHandler.NewPostApiHandler(cfg.APIv1, postSrv, cfg.Deps.AuthApiMdw)

	return &Module{
		PostService: postSrv,
	}, nil
}
