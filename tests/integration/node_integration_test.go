package integration

import (
	"context"
	"fmt"
	"github.com/PUMATeam/catapult/pkg/node"
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

var catapultNode *grpc.Server

func setup() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	catapultNode = grpc.NewServer()

	node.RegisterNodeServer(catapultNode, &MockNodeServer{})
	go func() {
		if err := catapultNode.Serve(lis); err != nil {
			log.Fatal("failed to start grpc server", err)
		}
	}()
}

func TestNodeApi(t *testing.T) {
	RegisterFailHandler(Fail)
	setup()
	RunSpecs(t, "Node Spec")
}

type MockNodeServer struct {
}

func (*MockNodeServer) StartVM(context.Context, *node.VmConfig) (*node.Response, error) {
	panic("implement me")
}

func (*MockNodeServer) StopVM(context.Context, *node.UUID) (*node.Response, error) {
	panic("implement me")
}

func (*MockNodeServer) ListVMs(context.Context, *empty.Empty) (*node.VmList, error) {
	panic("implement me")
}

func (*MockNodeServer) Health(context.Context, *empty.Empty) (*node.HealthCheckResponse, error) {
	panic("implement me")
}
