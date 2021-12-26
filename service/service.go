package service

import (
	"fmt"

	"github.com/ez-deploy/identity/db"
	pb "github.com/ez-deploy/protobuf/identity"
	"github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/public"
	"github.com/wuhuizuo/sqlm"
)

// Service impl protobuf.identity.Ops .
type Service struct {
	pb.UnimplementedIdentityOpsServer

	identityClient public.ClientService
	whoamiURL      string

	apiTokenTable *sqlm.Table
}

func New(kratosPublicHostname string, apiTokenTableDBDSN string) (*Service, error) {
	kratosCfg := &client.TransportConfig{
		Host:     kratosPublicHostname,
		BasePath: "/",
		Schemes:  []string{"http"},
	}
	kratosPublicClient := client.NewHTTPClientWithConfig(nil, kratosCfg).Public

	whoamiURL := fmt.Sprint("http://", kratosPublicHostname, "/sessions/whoami")

	database := &sqlm.Database{
		Driver: "mysql",
		DSN:    apiTokenTableDBDSN,
	}
	if err := database.Create(); err != nil {
		return nil, err
	}

	apiTokenTable := &sqlm.Table{
		Database:  database,
		TableName: "api_token",
	}
	apiTokenTable.SetRowModel(db.APITokenRawModel)

	resService := &Service{
		identityClient: kratosPublicClient,
		whoamiURL:      whoamiURL,
		apiTokenTable:  apiTokenTable,
	}

	return resService, nil
}
