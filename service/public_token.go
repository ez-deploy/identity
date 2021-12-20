package service

import (
	"context"

	pb "github.com/ez-deploy/protobuf/identity"
	"github.com/ez-deploy/protobuf/model"
)

// generate public_token.
func (s *Service) GeneratePublicToken(context.Context, *pb.GeneratePublicTokenReq) (*model.CommonResp, error) {
	return nil, nil
}

// list user's public_tokens.
func (s *Service) ListPublicToken(context.Context, *model.Identity) (*pb.ListPublicTokenResp, error) {
	return nil, nil
}

// delete public_token.
func (s *Service) DeletePublicToken(context.Context, *pb.DeletePublicTokenReq) (*model.CommonResp, error) {
	return nil, nil
}
