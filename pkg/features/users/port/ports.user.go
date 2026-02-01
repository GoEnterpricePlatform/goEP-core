package port

import (
	"context"
	"io"
	"time"

	"github.com/amorindev/go-tmpl/pkg/features/users/domain"
)

type UserRepo interface {
	Find(ctx context.Context, id string) (*domain.User, error)
	Insert(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Exists(ctx context.Context, id string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateAvatarPath(ctx context.Context, userID string, imgPath string, updatedAt time.Time) error
	ConfirmEmail(ctx context.Context, userID string) error
	InsertWithRoles(ctx context.Context, user *domain.User) error

	// Admin
	ExistsAdmin(ctx context.Context) (bool, error)
}

type UserSrv interface {
	UploadAvatar(ctx context.Context, path string, file io.Reader, userID string, contentType string) error
}

type UserFileStg interface {
	GetImage(ctx context.Context, imgPath string) (string, error)
	UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error
}
