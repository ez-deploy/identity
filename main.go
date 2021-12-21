package main

import (
	"io/ioutil"
	"log"
	"net"

	"github.com/ez-deploy/identity/service"
	pb "github.com/ez-deploy/protobuf/identity"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	_ "github.com/go-sql-driver/mysql"
)

const configFileName = "identityCfg.yaml"

type configBody struct {
	KratosPublicHostName string `yaml:"kratos_public_hostname,omitempty"`
	IdentityDSN          string `yaml:"identity_dsn,omitempty"`
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc, err := createServiceFromConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterOpsServer(s, svc)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func createServiceFromConfig(configFileName string) (*service.Service, error) {
	rawConfig, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	config := &configBody{}
	if err := yaml.Unmarshal(rawConfig, config); err != nil {
		return nil, err
	}

	return service.New(config.KratosPublicHostName, config.IdentityDSN)
}
