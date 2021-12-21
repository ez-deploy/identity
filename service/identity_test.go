package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ez-deploy/protobuf/identity"
	"github.com/ez-deploy/protobuf/model"
	"github.com/ory/kratos-client-go/client"
	"github.com/thanhpk/randstr"
)

var (
	testUsername string
	testEmail    string
	testPassword string

	verifyToken   model.Token
	serviceClient *Service
)

const randStrLen = 16

func TestMain(m *testing.M) {
	testUsername = randstr.Hex(randStrLen)
	testEmail = fmt.Sprint(randstr.Hex(randStrLen), "@test.com")
	testPassword = fmt.Sprintf(randstr.Hex(randStrLen), "_", randstr.Hex(randStrLen))

	cfg := &client.TransportConfig{
		Host:     "localhost:4433",
		BasePath: "/",
		Schemes:  []string{"http"},
	}
	serviceClient = &Service{
		identityClient: client.NewHTTPClientWithConfig(nil, cfg).Public,
		whoamiURL:      "http://localhost:4433/sessions/whoami",
	}

	m.Run()
}

func TestService_Register(t *testing.T) {
	t.Run("register", func(t *testing.T) {
		req := &identity.RegisterReq{
			Identity: &model.Identity{
				Email: testEmail,
				Name:  testUsername,
			},
			Password: testPassword,
		}

		resp, err := serviceClient.Register(context.Background(), req)
		if err != nil || resp.Error != nil {
			t.Fatal("register error,", err, resp.GetError())
		}
	})
}

func TestService_Login(t *testing.T) {
	t.Run("login", func(t *testing.T) {
		req := &identity.LoginReq{
			Email:    testEmail,
			Password: testPassword,
		}

		resp, err := serviceClient.Login(context.Background(), req)
		if err != nil || resp.Error != nil {
			t.Fatal("login error, ", err, resp.GetError())
		}

		if resp.Identity.Email != testEmail || resp.Identity.Name != testUsername {
			t.Fatal("login res neq", resp.Identity.Email, testEmail, resp.Identity.Name, testUsername)
		}
		verifyToken = *resp.Token
	})
}

func TestService_Verify(t *testing.T) {
	t.Run("login", func(t *testing.T) {
		req := &identity.VerifyReq{
			Token: &verifyToken,
		}

		resp, err := serviceClient.Verify(context.Background(), req)
		if err != nil || resp.Error != nil {
			t.Fatal("verify error, ", err, resp.GetError())
		}

		if resp.Identity.Email != testEmail || resp.Identity.Name != testUsername {
			t.Fatal("verify res neq", resp.Identity.Email, testEmail, resp.Identity.Name, testUsername)
		}
	})
}
