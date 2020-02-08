package storage

import (
	context "context"

	log "github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
	grpc "google.golang.org/grpc"
)

// TODO: move to internal?
type Service interface {
	Create(ctx context.Context, volume *Volume) (*Response, error)
	Delete(ctx context.Context, volID *uuid.UUID) (*Response, error)
	List(ctx context.Context) (*VolumeList, error)
}

type Storage struct {
	log *log.Logger
}

func NewStorageService(log *log.Logger) *Storage {
	return &Storage{log: log}
}

func (s *Storage) Create(ctx context.Context, volume *Volume) (*Response, error) {
	conn := createConn(ctx)
	client := NewStorageClient(conn)
	resp, err := client.Create(ctx, &Volume{UUID: volume.GetUUID(), Size: volume.GetSize()})
	if err != nil {
		return nil, err
	}

	s.log.WithContext(ctx).WithField("Response", resp).Info("Received response")

	return resp, nil
}

func (s *Storage) Delete(ctx context.Context, volID *uuid.UUID) (*Response, error) {
	return nil, nil
}

func (s *Storage) List(ctx context.Context) (*VolumeList, error) {
	return nil, nil
}

func createConn(ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx,
		"localhost:50051",
		grpc.WithInsecure())

	if err != nil {
		panic(err)
	}

	return conn
}
