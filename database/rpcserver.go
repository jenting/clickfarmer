// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package database

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"storj.io/clickfarmer/pb"
)

func Run(ctx context.Context, rpcAddr string) error {
	// create the in-memory database
	clickFarmer := &ClickFarmerDatabase{}

	// create a grpc server (without TLS)
	s := grpc.NewServer()

	// register the proto-specific methods on the grpc server
	pb.RegisterClickFarmerServer(s, clickFarmer)

	// listen on a tcp socket
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return err
	}

	// run the server
	return s.Serve(lis)
}

type ClickFarmerDatabase struct {
	pb.ClickFarmerServer

	mu                                 sync.Mutex
	redClicks, greenClicks, blueClicks int64
}

func (s *ClickFarmerDatabase) GetClicks(ctx context.Context, r *pb.GetClicksRequest) (*pb.GetClicksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &pb.GetClicksResponse{
		ClickCounts: &pb.ClickCounts{
			Red:   s.redClicks,
			Green: s.greenClicks,
			Blue:  s.blueClicks,
		},
	}, nil
}

func (s *ClickFarmerDatabase) SetClicks(ctx context.Context, r *pb.SetClicksRequest) (*pb.SetClicksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO
	return &pb.SetClicksResponse{}, nil
}
