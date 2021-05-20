package service

import (
	"context"
	"sync"

	"github.com/borud/t3/pkg/apipb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service is a very simple service implementation
type Service struct {
	mu      sync.RWMutex
	lastID  uint64
	entries map[uint64]*apipb.Map
}

func New() *Service {
	return &Service{
		entries: make(map[uint64]*apipb.Map),
	}
}

func (s *Service) AddMap(ctx context.Context, m *apipb.Map) (*apipb.AddMapResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	m.Id = s.lastID

	s.entries[s.lastID] = m

	return &apipb.AddMapResponse{Id: s.lastID}, nil
}

func (s *Service) GetMap(ctx context.Context, req *apipb.GetMapRequest) (*apipb.Map, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.entries[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return e, nil
}

func (s *Service) Update(ctx context.Context, m *apipb.Map) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check that the entry exists before we attempt to update
	existing, ok := s.entries[m.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	proto.Merge(existing, m)
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteMap(ctx context.Context, req *apipb.DeleteMapRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.entries, req.Id)

	return &emptypb.Empty{}, nil
}
