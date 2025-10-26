package service

import "context"

func (s *Service) Delete(ctx context.Context, id string) error {
    return s.OtpCodeRepo.Delete(ctx,id)
}