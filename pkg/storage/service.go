package storage

import (
	context "context"

	"github.com/PUMATeam/catapult/pkg/rpc"
	log "github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
)

// TODO: move to internal?
type Service interface {
	Create(ctx context.Context, volume *Volume) (*Response, error)
	Delete(ctx context.Context, volID *uuid.UUID) (*Response, error)
	List(ctx context.Context) (*VolumeList, error)
}

type Storage struct {
	log         *log.Logger
	connManager *rpc.GRPCConnection
}

func NewStorageService(connManager *rpc.GRPCConnection, log *log.Logger) *Storage {
	return &Storage{log: log, connManager: connManager}
}

func (s *Storage) Create(ctx context.Context, volume *Volume) (*Response, error) {
	conn, err := s.connManager.Connect(ctx, "localhost:5001")
	if err != nil {
		return nil, err
	}
	client := NewStorageClient(conn)
	resp, err := client.Create(ctx, &Volume{UUID: volume.GetUUID(), Size: volume.GetSize()})
	if err != nil {
		return nil, err
	}

	s.log.WithContext(ctx).WithField("Response", resp).Info("Received response")
	conn.Close()
	return resp, nil
}

func (s *Storage) Delete(ctx context.Context, volID *uuid.UUID) (*Response, error) {
	return nil, nil
}

func (s *Storage) List(ctx context.Context) (*VolumeList, error) {
	return nil, nil
}
