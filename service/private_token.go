package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ez-deploy/identity/db"
	"github.com/ez-deploy/protobuf/model"
	"github.com/thanhpk/randstr"
	"github.com/wuhuizuo/sqlm"
)

const tokenLen = 32

// ReGeneratePrivateToken regenerate private_token.
func (s *Service) ReGeneratePrivateToken(ctx context.Context, req *model.Identity) (*model.CommonResp, error) {
	oldToken, err := s.getPrivateTokenByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return s.generatePrivateToken(req.Email)
	}
	// err is not NotFound err, return it.
	if err != nil {
		return nil, err
	}

	return s.regeneratePrivateToken(oldToken.ID)
}

// GetPrivateToken get private_token by user's email.
func (s *Service) GetPrivateToken(ctx context.Context, req *model.Identity) (*model.Token, error) {
	privateToken, err := s.getPrivateTokenByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	resToken := &model.Token{
		Type:  model.TokenType_private,
		Token: privateToken.Token,
	}

	return resToken, nil
}

func (s *Service) getPrivateTokenByEmail(email string) (*db.APIToken, error) {
	filter := sqlm.SelectorFilter{
		"user_email": email,
		"type":       int32(model.TokenType_private),
	}

	oldToken := &db.APIToken{}
	err := s.apiTokenTable.Get(filter, oldToken)
	if err != nil {
		return nil, err
	}

	return oldToken, nil
}

func (s *Service) generatePrivateToken(email string) (*model.CommonResp, error) {
	newAPIToken := db.APIToken{
		UserEmail: email,
		Type:      int32(model.TokenType_private),
		Token:     randstr.Hex(tokenLen),
	}

	if _, err := s.apiTokenTable.Insert(newAPIToken); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}

func (s *Service) regeneratePrivateToken(oldTokenID int64) (*model.CommonResp, error) {
	filter := sqlm.SelectorFilter{"id": oldTokenID}
	updateParts := map[string]interface{}{"token": randstr.Hex(tokenLen)}

	if err := s.apiTokenTable.Update(filter, updateParts); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}
