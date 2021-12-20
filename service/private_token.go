package service

import (
	"context"

	"github.com/ez-deploy/protobuf/model"
)

// generate private_token.
func (s *Service) ReGeneratePrivateToken(context.Context, *model.Identity) (*model.CommonResp, error) {
	return nil, nil
}

// get private_token.
func (s *Service) GetPrivateToken(context.Context, *model.Identity) (*model.Token, error) {
	return nil, nil
}
