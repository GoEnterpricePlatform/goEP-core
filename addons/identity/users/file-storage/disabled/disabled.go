package disabled

import (
	"context"
	"errors"
	"io"
)

type DisabledFS struct{}

func NewDisabledAdapter() *DisabledFS {
	return &DisabledFS{}
}

func (a *DisabledFS) GetImage(ctx context.Context, imgPath string) (string, error) {
	return "", errors.New("error user file storage disabled")
}

func (a *DisabledFS) UploadImage(ctx context.Context, imgPath string, file io.Reader, contentType string) error {
	return errors.New("error user file storage disabled")
}
