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

func TestNodes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Node Suite")
}

var _ = BeforeSuite(func() {
	setup()
})

var _ = AfterSuite(func() {
	catapultNode.Stop()
})

func setup() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8888))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	catapultNode = grpc.NewServer()

	node.RegisterNodeServer(catapultNode, &MockNodeServer{})

	go func() {
		err = catapultNode.Serve(lis)
		if err != nil {
			fmt.Errorf("failed to resiter node server %s", err)
		}
	}()
}

var _ = Describe("Node", func() {
		var (
			client node.NodeClient
		)
		// restart a new service each time?
		conn, e := grpc.Dial("localhost:8888", grpc.WithInsecure())
		if e != nil {
			Fail(e.Error())
		}
		client = node.NewNodeClient(conn)


		Describe("API sanity", func() {
			Context("Endpoints", func() {
				It("responds to health checks", func() {
					response, e := client.Health(context.TODO(), &empty.Empty{})
					Expect(e).NotTo(HaveOccurred())
					Expect(response.State).To(Equal(node.HealthCheckState_UP))
				})

				It("responds to Start VM ", func() {
					response, e := client.StartVM(context.TODO(), &node.VmConfig{})
					Expect(e).NotTo(HaveOccurred())
					Expect(response.Status).To(Equal(node.RequestStatus_SUCCESSFUL))
				})

				It("responds to List ", func() {
					list, e := client.ListVMs(context.TODO(), &empty.Empty{})
					Expect(e).NotTo(HaveOccurred())
					Expect(list.VmID).NotTo(BeNil())
				})
			})
		})
})

type MockNodeServer struct {

}

func (*MockNodeServer) StartVM(context.Context, *node.VmConfig) (*node.Response, error) {
	return &node.Response{Status: node.RequestStatus_SUCCESSFUL}, nil
}

func (*MockNodeServer) StopVM(context.Context, *node.UUID) (*node.Response, error) {
	return &node.Response{Status: node.RequestStatus_SUCCESSFUL}, nil
}

func (*MockNodeServer) ListVMs(context.Context, *empty.Empty) (*node.VmList, error) {
	list := node.VmList{}
	list.VmID = &node.UUID{}
	return &list, nil
}

func (*MockNodeServer) Health(context.Context, *empty.Empty) (*node.HealthCheckResponse, error) {
	return &node.HealthCheckResponse{State: node.HealthCheckState_UP}, nil
}

