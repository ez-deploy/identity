package service

import (
	"github.com/ory/kratos-client-go/client/public"
	"github.com/wuhuizuo/sqlm"
)

// Service impl protobuf.identity.Ops .
type Service struct {
	identityClient  public.ClientService
	kratosPublicURL string

	apiTokenTable *sqlm.Table
}
