package service

import "github.com/ory/kratos-client-go/client/public"

// Service impl protobuf.identity.Ops .
type Service struct {
	identityClient  public.ClientService
	kratosPublicURL string
}

// NewService with kratos client.
func New(identityClient public.ClientService) *Service {
	return &Service{identityClient: identityClient}
}
