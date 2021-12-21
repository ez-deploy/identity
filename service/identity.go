package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ez-deploy/protobuf/convert"
	pb "github.com/ez-deploy/protobuf/identity"
	"github.com/ez-deploy/protobuf/model"
	"github.com/ory/kratos-client-go/client/public"
	"github.com/ory/kratos-client-go/models"
	"github.com/pkg/errors"
)

// Register by email and password.
func (s *Service) Register(ctx context.Context, req *pb.RegisterReq) (*model.CommonResp, error) {
	if req.Identity.Email == "" {
		return model.NewCommonRespWithErrorMessage("email is required"), nil
	}

	initParams := &public.InitializeSelfServiceRegistrationViaAPIFlowParams{Context: ctx}

	flow, err := s.identityClient.InitializeSelfServiceRegistrationViaAPIFlow(initParams)
	if err != nil {
		return nil, err
	}
	flowID := (*string)(flow.Payload.ID)

	registerPayload := map[string]string{
		"traits.email": req.Identity.Email,
		"traits.name":  req.Identity.Name,
		"password":     req.Password,
	}
	registerParams := &public.CompleteSelfServiceRegistrationFlowWithPasswordMethodParams{
		Payload: registerPayload,
		Flow:    flowID,
		Context: ctx,
	}

	if _, err := s.identityClient.CompleteSelfServiceRegistrationFlowWithPasswordMethod(registerParams); err != nil {
		return model.NewCommonRespWithError(err), nil
	}

	return &model.CommonResp{}, nil
}

// Login by email and password.
func (s *Service) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	initParams := &public.InitializeSelfServiceLoginViaAPIFlowParams{Context: ctx}

	flow, err := s.identityClient.InitializeSelfServiceLoginViaAPIFlow(initParams)
	if err != nil {
		return nil, err
	}
	flowID := string(*flow.Payload.ID)

	loginBody := &models.CompleteSelfServiceLoginFlowWithPasswordMethod{
		Identifier: req.Email,
		Password:   req.Password,
	}
	loginParam := &public.CompleteSelfServiceLoginFlowWithPasswordMethodParams{
		Body:    loginBody,
		Flow:    flowID,
		Context: ctx,
	}

	loginRes, err := s.identityClient.CompleteSelfServiceLoginFlowWithPasswordMethod(loginParam)
	if err != nil {
		return &pb.LoginResp{Error: model.NewError(err)}, nil
	}

	return newLoginResp(loginRes)
}

// Verify by session_token.
func (s *Service) Verify(ctx context.Context, req *pb.VerifyReq) (*pb.VerifyResp, error) {
	switch req.Token.Type {
	case model.TokenType_session:
		return s.verifySessionToken(req.Token.Token)
	default:
		return nil, errors.New("not impl")
	}
}

func newLoginResp(loginRes *public.CompleteSelfServiceLoginFlowWithPasswordMethodOK) (*pb.LoginResp, error) {
	resIdentity := &model.Identity{}
	if err := convert.WithJSON(loginRes.Payload.Session.Identity.Traits, resIdentity); err != nil {
		return nil, errors.WithMessage(err, "convert traits info")
	}

	resToken := &model.Token{
		Type:  model.TokenType_session,
		Token: *loginRes.Payload.SessionToken,
	}

	resp := &pb.LoginResp{
		Identity: resIdentity,
		Token:    resToken,
	}

	return resp, nil
}

func (s *Service) verifySessionToken(rawToken string) (*pb.VerifyResp, error) {
	req, err := http.NewRequest("GET", s.whoamiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Session-Token", rawToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &pb.VerifyResp{Error: model.NewErrorWithMessage("unauthorized")}, nil
	}

	resIdentity, err := loadIdentityFromResponse(resp)
	if err != nil {
		return nil, err
	}

	res := &pb.VerifyResp{
		Identity:  resIdentity,
		TokenType: model.TokenType_session,
	}

	return res, nil
}

type verifyResposne struct {
	Identity struct {
		Traits model.Identity `json:"traits"`
	} `json:"identity"`
}

func loadIdentityFromResponse(resp *http.Response) (*model.Identity, error) {
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	verifyResp := &verifyResposne{}
	if err := json.Unmarshal(rawBody, verifyResp); err != nil {
		return nil, err
	}

	return &verifyResp.Identity.Traits, nil
}
