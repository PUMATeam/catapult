package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// GRPCConnection manages all grpc connections
type GRPCConnection struct {
	// Revisit this idea at some point,
	// on paper this seems like a good fit
	// but maybe there's something I am missing
	connections sync.Map
}

// NewGRPCConnection creates a new GRPCConnection instance
func NewGRPCConnection() *GRPCConnection {
	return &GRPCConnection{}
}

// Connect connects to an `address` and stores it in the internal map
func (g *GRPCConnection) Connect(ctx context.Context, address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		address,
		grpc.WithBlock(),
		grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	if conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("Could not connect to node in allocated time")
	}

	g.connections.Store(address, conn)

	return conn, nil
}

// GetConnection receives a address and returns its connection.
// If a connection is not available return nil
func (g *GRPCConnection) GetConnection(address string) *grpc.ClientConn {
	v, ok := g.connections.Load(address)
	if !ok {
		return nil
	}

	return v.(*grpc.ClientConn)
}

// Close recieves a connection and closes it
func (g *GRPCConnection) Close(ctx context.Context, address string) error {
	conn := g.GetConnection(address)
	if conn != nil {
		err := conn.Close()
		if err != nil {
			return err
		}
	}

	g.connections.Delete(address)

	return nil
}

// Shutdown closes all the connections
func (g *GRPCConnection) Shutdown() []error {
	errors := make([]error, 0, 0)
	g.connections.Range(func(key, value interface{}) bool {
		err := value.(*grpc.ClientConn).Close()

		if err != nil {
			errors = append(errors, err)
		}

		return true
	})

	return errors
}
