package service

import (
	"context"

	"github.com/ez-deploy/identity/db"
	"github.com/ez-deploy/protobuf/convert"
	pb "github.com/ez-deploy/protobuf/identity"
	"github.com/ez-deploy/protobuf/model"
	"github.com/thanhpk/randstr"
	"github.com/wuhuizuo/sqlm"
)

// generate public_token.
func (s *Service) GeneratePublicToken(ctx context.Context, req *pb.GeneratePublicTokenReq) (*model.CommonResp, error) {
	newDBToken := newDBAPITokenFromPBAPIToken(req.ApiToken)

	newDBToken.ID = 0
	newDBToken.Token = randstr.Hex(tokenLen)
	newDBToken.UserEmail = req.Identity.Email

	if _, err := s.apiTokenTable.Insert(newDBToken); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}

// list user's public_tokens.
func (s *Service) ListPublicToken(ctx context.Context, req *model.Identity) (*pb.ListPublicTokenResp, error) {
	filter := sqlm.SelectorFilter{"user_email": req.Email}
	listOptions := sqlm.ListOptions{
		OrderByColumn: "deadline_timestamp",
		AllColumns:    true,
	}

	records, err := s.apiTokenTable.List(filter, listOptions)
	if err != nil {
		return nil, err
	}

	resTokens := []*model.APIToken{}
	for _, record := range records {
		resToken := &db.APIToken{}
		if err := convert.WithJSON(record, resToken); err != nil {
			return nil, err
		}

		resTokens = append(resTokens, newPBAPITokenFromDBAPIToken(resToken))
	}

	return &pb.ListPublicTokenResp{PublicTokens: resTokens}, nil
}

// delete public_token.
func (s *Service) DeletePublicToken(ctx context.Context, req *pb.DeletePublicTokenReq) (*model.CommonResp, error) {
	filter := sqlm.SelectorFilter{"id": req.TokenId}

	if err := s.apiTokenTable.Delete(filter); err != nil {
		return nil, err
	}

	return &model.CommonResp{}, nil
}

func newDBAPITokenFromPBAPIToken(pbAPIToken *model.APIToken) *db.APIToken {
	allowedActions := int32(0)
	for _, allowedAction := range pbAPIToken.AllowedActions {
		allowedActions = (allowedActions | int32(allowedAction))
	}

	resDBAPIToken := &db.APIToken{
		ID:    pbAPIToken.Id,
		Type:  int32(pbAPIToken.Token.Type),
		Token: pbAPIToken.Token.Token,

		Name:              pbAPIToken.Name,
		Message:           pbAPIToken.Message,
		DeadlineTimestamp: pbAPIToken.DeadlineTimestamp,
		AllowedActions:    allowedActions,
	}

	return resDBAPIToken
}

func newPBAPITokenFromDBAPIToken(dbAPIToken *db.APIToken) *model.APIToken {
	allowedActions := []model.Actions{}
	for action := range model.Actions_name {
		if (dbAPIToken.AllowedActions & action) != 0 {
			allowedActions = append(allowedActions, model.Actions(action))
		}
	}

	resToken := &model.Token{
		Type:  model.TokenType(dbAPIToken.Type),
		Token: dbAPIToken.Token,
	}

	resPBToken := &model.APIToken{
		Id:    dbAPIToken.ID,
		Token: resToken,

		Name:              dbAPIToken.Name,
		Message:           dbAPIToken.Message,
		DeadlineTimestamp: dbAPIToken.DeadlineTimestamp,
		AllowedActions:    allowedActions,
	}

	return resPBToken
}
